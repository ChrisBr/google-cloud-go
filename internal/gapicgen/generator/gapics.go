// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generator

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/internal/gapicgen/execv"
	"cloud.google.com/go/internal/gapicgen/execv/gocmd"
	"cloud.google.com/go/internal/gapicgen/gensnippets"
	"cloud.google.com/go/internal/gapicgen/git"
	"gopkg.in/yaml.v2"
)

var (
	//go:embed _CHANGES.md.txt
	changesTmpl string
	//go:embed _README.md.txt
	readmeTmpl string
	//go:embed _version.go.txt
	versionTmpl string
)

// GapicGenerator is used to regenerate gapic libraries.
type GapicGenerator struct {
	googleapisDir     string
	protoDir          string
	googleCloudDir    string
	genprotoDir       string
	gapicToGenerate   string
	regenOnly         bool
	onlyGenerateGapic bool
	genModule         bool
	modifiedPkgs      []string
	forceAll          bool
}

// NewGapicGenerator creates a GapicGenerator.
func NewGapicGenerator(c *Config, modifiedPkgs []string) *GapicGenerator {
	return &GapicGenerator{
		googleapisDir:     c.GoogleapisDir,
		protoDir:          c.ProtoDir,
		googleCloudDir:    c.GapicDir,
		genprotoDir:       c.GenprotoDir,
		gapicToGenerate:   c.GapicToGenerate,
		regenOnly:         c.RegenOnly,
		onlyGenerateGapic: c.OnlyGenerateGapic,
		genModule:         c.GenModule,
		modifiedPkgs:      modifiedPkgs,
		forceAll:          c.ForceAll,
	}
}

type modInfo struct {
	path              string
	importPath        string
	serviceImportPath string
}

// Regen generates gapics.
func (g *GapicGenerator) Regen(ctx context.Context) error {
	log.Println("regenerating gapics")
	var newMods []modInfo
	for _, c := range microgenGapicConfigs {
		if !g.shouldGenerateConfig(c) {
			continue
		}

		modImportPath := filepath.Join("cloud.google.com/go", strings.Split(strings.TrimPrefix(c.importPath, "cloud.google.com/go/"), "/")[0])
		modPath := filepath.Join(g.googleCloudDir, modImportPath)
		if g.genModule {
			if err := generateModule(modPath, modImportPath); err != nil {
				return err
			}
			newMods = append(newMods, modInfo{
				path:              filepath.Join(g.googleCloudDir, strings.TrimPrefix(modImportPath, "cloud.google.com/go")),
				importPath:        modImportPath,
				serviceImportPath: c.importPath,
			})
		}
		if err := g.microgen(c); err != nil {
			return err
		}
		if g.genModule {
			if err := gocmd.ModTidy(modPath); err != nil {
				return nil
			}
		}
		if err := g.genVersionFile(c); err != nil {
			return err
		}
	}

	if err := g.copyMicrogenFiles(); err != nil {
		return err
	}

	// Get rid of diffs related to bad formatting.
	if err := gocmd.Vet(g.googleCloudDir); err != nil {
		return err
	}

	if g.regenOnly {
		return nil
	}

	manifest, err := g.manifest(microgenGapicConfigs)
	if err != nil {
		return err
	}

	if g.genModule {
		for _, modInfo := range newMods {
			generateReadmeAndChanges(modInfo.path, modInfo.importPath, manifest[modInfo.serviceImportPath].Description)
		}
	}

	if !g.onlyGenerateGapic {
		if err := g.regenSnippets(ctx); err != nil {
			return err
		}
	}

	if !g.onlyGenerateGapic {
		if err := execv.ForEachMod(g.googleCloudDir, g.addModReplaceGenproto); err != nil {
			return err
		}
	}

	if err := gocmd.Vet(g.googleCloudDir); err != nil {
		return err
	}

	if err := gocmd.Build(g.googleCloudDir); err != nil {
		return err
	}

	if !g.onlyGenerateGapic {
		if err := execv.ForEachMod(g.googleCloudDir, g.dropModReplaceGenproto); err != nil {
			return err
		}
	}

	return nil
}

func (g *GapicGenerator) shouldGenerateConfig(c *microgenConfig) bool {
	if g.forceAll && !c.stopGeneration {
		return true
	}

	// Skip generation if generating all of the gapics and the associated
	// config has a block on it. Or if generating a single gapic and it does
	// not match the specified import path.
	if (c.stopGeneration && g.gapicToGenerate == "") ||
		(g.gapicToGenerate != "" && !strings.Contains(g.gapicToGenerate, c.importPath)) ||
		(g.forceAll && !c.stopGeneration) {
		return false
	}
	return true
}

// RegenSnippets regenerates the snippets for all GAPICs configured to be generated.
func (g *GapicGenerator) regenSnippets(ctx context.Context) error {
	log.Println("regenerating snippets")

	snippetDir := filepath.Join(g.googleCloudDir, "internal", "generated", "snippets")
	apiShortnames, err := g.parseAPIShortnames(microgenGapicConfigs, manualEntries)
	if err != nil {
		return err
	}
	if err := gensnippets.Generate(g.googleCloudDir, snippetDir, apiShortnames); err != nil {
		log.Printf("warning: got the following non-fatal errors generating snippets: %v", err)
	}
	if err := replaceAllForSnippets(g.googleCloudDir, snippetDir); err != nil {
		return err
	}
	if err := gocmd.ModTidy(snippetDir); err != nil {
		return err
	}
	return nil
}

func replaceAllForSnippets(googleCloudDir, snippetDir string) error {
	return execv.ForEachMod(googleCloudDir, func(dir string) error {
		if dir == snippetDir {
			return nil
		}

		mod, err := gocmd.ListModName(dir)
		if err != nil {
			return err
		}

		// Replace it. Use a relative path to avoid issues on different systems.
		rel, err := filepath.Rel(snippetDir, dir)
		if err != nil {
			return err
		}
		c := execv.Command("bash", "-c", `go mod edit -replace "$MODULE=$MODULE_PATH"`)
		c.Dir = snippetDir
		c.Env = []string{
			fmt.Sprintf("PATH=%s", os.Getenv("PATH")), // TODO(deklerk): Why do we need to do this? Doesn't seem to be necessary in other exec.Commands.
			fmt.Sprintf("HOME=%s", os.Getenv("HOME")), // TODO(deklerk): Why do we need to do this? Doesn't seem to be necessary in other exec.Commands.
			fmt.Sprintf("MODULE=%s", mod),
			fmt.Sprintf("MODULE_PATH=%s", rel),
		}
		return c.Run()
	})
}

// addModReplaceGenproto adds a genproto replace statement that points genproto
// to the local copy. This is necessary since the remote genproto may not have
// changes that are necessary for the in-flight regen.
func (g *GapicGenerator) addModReplaceGenproto(dir string) error {
	log.Printf("[%s] adding temporary genproto replace statement", dir)
	c := execv.Command("bash", "-c", `
set -ex

go mod edit -replace "google.golang.org/genproto=$GENPROTO_DIR"
`)
	c.Dir = dir
	c.Env = []string{
		"GENPROTO_DIR=" + g.genprotoDir,
		fmt.Sprintf("PATH=%s", os.Getenv("PATH")), // TODO(deklerk): Why do we need to do this? Doesn't seem to be necessary in other exec.Commands.
		fmt.Sprintf("HOME=%s", os.Getenv("HOME")), // TODO(deklerk): Why do we need to do this? Doesn't seem to be necessary in other exec.Commands.
	}
	return c.Run()
}

// dropModReplaceGenproto drops the genproto replace statement. It is intended
// to be run after addModReplaceGenproto.
func (g *GapicGenerator) dropModReplaceGenproto(dir string) error {
	log.Printf("[%s] removing genproto replace statement", dir)
	c := execv.Command("bash", "-c", `
set -ex

go mod edit -dropreplace "google.golang.org/genproto"
`)
	c.Dir = dir
	c.Env = []string{
		fmt.Sprintf("PATH=%s", os.Getenv("PATH")), // TODO(deklerk): Why do we need to do this? Doesn't seem to be necessary in other exec.Commands.
		fmt.Sprintf("HOME=%s", os.Getenv("HOME")), // TODO(deklerk): Why do we need to do this? Doesn't seem to be necessary in other exec.Commands.
	}
	return c.Run()
}

// microgen runs the microgenerator on a single microgen config.
func (g *GapicGenerator) microgen(conf *microgenConfig) error {
	log.Println("microgen generating", conf.pkg)

	var protoFiles []string
	if err := filepath.Walk(g.googleapisDir+"/"+conf.inputDirectoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Ignore compute_small.proto which is just for testing and would cause a collision if used in generation.
		//
		// TODO(noahdietz): Remove this when it is no longer needed.
		if strings.Contains(info.Name(), ".proto") && !strings.Contains(info.Name(), "compute_small.proto") {
			protoFiles = append(protoFiles, path)
		}
		return nil
	}); err != nil {
		return err
	}

	args := []string{"-I", g.googleapisDir,
		"--experimental_allow_proto3_optional",
		"-I", g.protoDir,
		"--go_gapic_out", g.googleCloudDir,
		"--go_gapic_opt", fmt.Sprintf("go-gapic-package=%s;%s", conf.importPath, conf.pkg),
		"--go_gapic_opt", fmt.Sprintf("api-service-config=%s", filepath.Join(conf.inputDirectoryPath, conf.apiServiceConfigPath))}

	if conf.releaseLevel != "" {
		args = append(args, "--go_gapic_opt", fmt.Sprintf("release-level=%s", conf.releaseLevel))
	}
	if conf.gRPCServiceConfigPath != "" {
		args = append(args, "--go_gapic_opt", fmt.Sprintf("grpc-service-config=%s", filepath.Join(conf.inputDirectoryPath, conf.gRPCServiceConfigPath)))
	}
	if !conf.disableMetadata {
		args = append(args, "--go_gapic_opt", "metadata")
	}
	if len(conf.transports) > 0 {
		args = append(args, "--go_gapic_opt", fmt.Sprintf("transport=%s", strings.Join(conf.transports, "+")))
	}
	// This is a bummer way of toggling diregapic generation, but it compute is the only one for the near term.
	if conf.pkg == "compute" {
		args = append(args, "--go_gapic_opt", "diregapic")
	}
	args = append(args, protoFiles...)
	c := execv.Command("protoc", args...)
	c.Dir = g.googleapisDir
	return c.Run()
}

func (g *GapicGenerator) genVersionFile(conf *microgenConfig) error {
	// These directories are not modules on purpose, don't generate a version
	// file for them.
	if conf.importPath == "cloud.google.com/go/longrunning/autogen" ||
		conf.importPath == "cloud.google.com/go/debugger/apiv2" {
		return nil
	}
	relDir := strings.TrimPrefix(conf.importPath, "cloud.google.com/go/")
	rootPackage := strings.Split(relDir, "/")[0]
	rootModInternal := fmt.Sprintf("cloud.google.com/go/%s/internal", rootPackage)

	f, err := os.Create(filepath.Join(g.googleCloudDir, relDir, "version.go"))
	if err != nil {
		return err
	}
	defer f.Close()

	t := template.Must(template.New("version").Parse(versionTmpl))
	versionData := struct {
		Year               int
		Package            string
		ModuleRootInternal string
	}{
		Year:               time.Now().Year(),
		Package:            conf.pkg,
		ModuleRootInternal: rootModInternal,
	}
	if err := t.Execute(f, versionData); err != nil {
		return err
	}
	return nil
}

// manifestEntry is used for JSON marshaling in manifest.
type manifestEntry struct {
	DistributionName  string      `json:"distribution_name"`
	Description       string      `json:"description"`
	Language          string      `json:"language"`
	ClientLibraryType string      `json:"client_library_type"`
	DocsURL           string      `json:"docs_url"`
	ReleaseLevel      string      `json:"release_level"`
	LibraryType       LibraryType `json:"library_type"`
}

type LibraryType string

const (
	GapicAutoLibraryType   LibraryType = "GAPIC_AUTO"
	GapicManualLibraryType LibraryType = "GAPIC_MANUAL"
	CoreLibraryType        LibraryType = "CORE"
	AgentLibraryType       LibraryType = "AGENT"
	OtherLibraryType       LibraryType = "OTHER"
)

// TODO: consider getting Description from the gapic, if there is one.
var manualEntries = []manifestEntry{
	// Pure manual clients.
	{
		DistributionName:  "cloud.google.com/go/bigquery",
		Description:       "BigQuery",
		Language:          "Go",
		ClientLibraryType: "manual",
		DocsURL:           "https://cloud.google.com/go/docs/reference/cloud.google.com/go/bigquery/latest",
		ReleaseLevel:      "ga",
		LibraryType:       GapicManualLibraryType,
	},
	{
		DistributionName:  "cloud.google.com/go/bigtable",
		Description:       "Cloud BigTable",
		Language:          "Go",
		ClientLibraryType: "manual",
		DocsURL:           "https://cloud.google.com/go/docs/reference/cloud.google.com/go/bigtable/latest",
		ReleaseLevel:      "ga",
		LibraryType:       GapicManualLibraryType,
	},
	{
		DistributionName:  "cloud.google.com/go/datastore",
		Description:       "Cloud Datastore",
		Language:          "Go",
		ClientLibraryType: "manual",
		DocsURL:           "https://cloud.google.com/go/docs/reference/cloud.google.com/go/datastore/latest",
		ReleaseLevel:      "ga",
		LibraryType:       GapicManualLibraryType,
	},
	{
		DistributionName:  "cloud.google.com/go/iam",
		Description:       "Cloud IAM",
		Language:          "Go",
		ClientLibraryType: "manual",
		DocsURL:           "https://cloud.google.com/go/docs/reference/cloud.google.com/go/iam/latest",
		ReleaseLevel:      "ga",
		LibraryType:       CoreLibraryType,
	},
	{
		DistributionName:  "cloud.google.com/go/storage",
		Description:       "Cloud Storage (GCS)",
		Language:          "Go",
		ClientLibraryType: "manual",
		DocsURL:           "https://cloud.google.com/go/docs/reference/cloud.google.com/go/storage/latest",
		ReleaseLevel:      "ga",
		LibraryType:       GapicManualLibraryType,
	},
	{
		DistributionName:  "cloud.google.com/go/rpcreplay",
		Description:       "RPC Replay",
		Language:          "Go",
		ClientLibraryType: "manual",
		DocsURL:           "https://cloud.google.com/go/docs/reference/cloud.google.com/go/latest/rpcreplay",
		ReleaseLevel:      "ga",
		LibraryType:       OtherLibraryType,
	},
	{
		DistributionName:  "cloud.google.com/go/profiler",
		Description:       "Cloud Profiler",
		Language:          "Go",
		ClientLibraryType: "manual",
		DocsURL:           "https://cloud.google.com/go/docs/reference/cloud.google.com/go/profiler/latest",
		ReleaseLevel:      "ga",
		LibraryType:       AgentLibraryType,
	},
	{
		DistributionName:  "cloud.google.com/go/compute/metadata",
		Description:       "Service Metadata API",
		Language:          "Go",
		ClientLibraryType: "manual",
		DocsURL:           "https://cloud.google.com/go/docs/reference/cloud.google.com/go/compute/latest/metadata",
		ReleaseLevel:      "ga",
		LibraryType:       CoreLibraryType,
	},
	{
		DistributionName:  "cloud.google.com/go/functions/metadata",
		Description:       "Cloud Functions",
		Language:          "Go",
		ClientLibraryType: "manual",
		DocsURL:           "https://cloud.google.com/go/docs/reference/cloud.google.com/go/functions/latest/metadata",
		ReleaseLevel:      "alpha",
		LibraryType:       CoreLibraryType,
	},
	// Manuals with a GAPIC.
	{
		DistributionName:  "cloud.google.com/go/errorreporting",
		Description:       "Cloud Error Reporting API",
		Language:          "Go",
		ClientLibraryType: "manual",
		DocsURL:           "https://cloud.google.com/go/docs/reference/cloud.google.com/go/errorreporting/latest",
		ReleaseLevel:      "beta",
		LibraryType:       GapicManualLibraryType,
	},
	{
		DistributionName:  "cloud.google.com/go/firestore",
		Description:       "Cloud Firestore API",
		Language:          "Go",
		ClientLibraryType: "manual",
		DocsURL:           "https://cloud.google.com/go/docs/reference/cloud.google.com/go/firestore/latest",
		ReleaseLevel:      "ga",
		LibraryType:       GapicManualLibraryType,
	},
	{
		DistributionName:  "cloud.google.com/go/logging",
		Description:       "Cloud Logging API",
		Language:          "Go",
		ClientLibraryType: "manual",
		DocsURL:           "https://cloud.google.com/go/docs/reference/cloud.google.com/go/logging/latest",
		ReleaseLevel:      "ga",
		LibraryType:       GapicManualLibraryType,
	},
	{
		DistributionName:  "cloud.google.com/go/pubsub",
		Description:       "Cloud PubSub",
		Language:          "Go",
		ClientLibraryType: "manual",
		DocsURL:           "https://cloud.google.com/go/docs/reference/cloud.google.com/go/pubsub/latest",
		ReleaseLevel:      "ga",
		LibraryType:       GapicManualLibraryType,
	},
	{
		DistributionName:  "cloud.google.com/go/spanner",
		Description:       "Cloud Spanner",
		Language:          "Go",
		ClientLibraryType: "manual",
		DocsURL:           "https://cloud.google.com/go/docs/reference/cloud.google.com/go/spanner/latest",
		ReleaseLevel:      "ga",
		LibraryType:       GapicManualLibraryType,
	},
	{
		DistributionName:  "cloud.google.com/go/pubsublite",
		Description:       "Cloud PubSub Lite",
		Language:          "Go",
		ClientLibraryType: "manual",
		DocsURL:           "https://cloud.google.com/go/docs/reference/cloud.google.com/go/pubsublite/latest",
		ReleaseLevel:      "beta",
		LibraryType:       GapicManualLibraryType,
	},
}

// manifest writes a manifest file with info about all of the confs.
func (g *GapicGenerator) manifest(confs []*microgenConfig) (map[string]manifestEntry, error) {
	log.Println("updating gapic manifest")
	entries := map[string]manifestEntry{} // Key is the package name.
	f, err := os.Create(filepath.Join(g.googleCloudDir, "internal", ".repo-metadata-full.json"))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	for _, manual := range manualEntries {
		entries[manual.DistributionName] = manual
	}
	for _, conf := range confs {
		yamlPath := filepath.Join(g.googleapisDir, conf.inputDirectoryPath, conf.apiServiceConfigPath)
		yamlFile, err := os.Open(yamlPath)
		if err != nil {
			return nil, err
		}
		yamlConfig := struct {
			Title string `yaml:"title"` // We only need the title field.
		}{}
		if err := yaml.NewDecoder(yamlFile).Decode(&yamlConfig); err != nil {
			return nil, fmt.Errorf("decode: %v", err)
		}
		docURL, err := docURL(g.googleCloudDir, conf.importPath)
		if err != nil {
			return nil, fmt.Errorf("unable to build docs URL: %v", err)
		}
		entry := manifestEntry{
			DistributionName:  conf.importPath,
			Description:       yamlConfig.Title,
			Language:          "Go",
			ClientLibraryType: "generated",
			DocsURL:           docURL,
			ReleaseLevel:      conf.releaseLevel,
		}
		entries[conf.importPath] = entry
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return entries, enc.Encode(entries)
}

// copyMicrogenFiles takes microgen files from gocloudDir/cloud.google.com/go
// and places them in gocloudDir.
func (g *GapicGenerator) copyMicrogenFiles() error {
	// The period at the end is analagous to * (copy everything in this dir).
	c := execv.Command("cp", "-R", g.googleCloudDir+"/cloud.google.com/go/.", ".")
	c.Dir = g.googleCloudDir
	if err := c.Run(); err != nil {
		return err
	}

	c = execv.Command("rm", "-rf", "cloud.google.com")
	c.Dir = g.googleCloudDir
	return c.Run()
}

func (g *GapicGenerator) parseAPIShortnames(confs []*microgenConfig, manualEntries []manifestEntry) (map[string]string, error) {
	shortnames := map[string]string{}
	for _, conf := range confs {
		yamlPath := filepath.Join(g.googleapisDir, conf.inputDirectoryPath, conf.apiServiceConfigPath)
		yamlFile, err := os.Open(yamlPath)
		if err != nil {
			return nil, err
		}
		config := struct {
			Name string `yaml:"name"`
		}{}
		if err := yaml.NewDecoder(yamlFile).Decode(&config); err != nil {
			return nil, fmt.Errorf("decode: %v", err)
		}
		shortname := strings.TrimSuffix(config.Name, ".googleapis.com")
		shortnames[conf.importPath] = shortname
	}

	// Do our best for manuals.
	for _, manual := range manualEntries {
		p := strings.TrimPrefix(manual.DistributionName, "cloud.google.com/go/")
		if strings.Contains(p, "/") {
			p = p[0:strings.Index(p, "/")]
		}
		shortnames[manual.DistributionName] = p
	}
	return shortnames, nil
}

func (g *GapicGenerator) findModifiedDirs() ([]string, error) {
	log.Println("finding modifiled directories")
	files, err := git.FindModifiedAndUntrackedFiles(g.googleCloudDir)
	if err != nil {
		return nil, err
	}
	dirs := map[string]bool{}
	for _, file := range files {
		dir := filepath.Dir(filepath.Join(g.googleCloudDir, file))
		dirs[dir] = true
	}

	// Add modified dirs from genproto. Sometimes only a request struct will be
	// updated, in these cases we should still make modifications the
	// corresponding gapic directories.
	for _, pkg := range g.modifiedPkgs {
		dir := filepath.Join(g.googleCloudDir, pkg)
		dirs[dir] = true
	}

	var dirList []string
	for dir := range dirs {
		dirList = append(dirList, dir)
	}
	return dirList, nil
}

func docURL(cloudDir, importPath string) (string, error) {
	suffix := strings.TrimPrefix(importPath, "cloud.google.com/go/")
	mod, err := gocmd.CurrentMod(filepath.Join(cloudDir, suffix))
	if err != nil {
		return "", err
	}
	pkgPath := strings.TrimPrefix(strings.TrimPrefix(importPath, mod), "/")
	return "https://cloud.google.com/go/docs/reference/" + mod + "/latest/" + pkgPath, nil
}

func generateModule(path, importPath string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return nil
	}
	log.Printf("Creating %s/go.mod", path)
	return gocmd.ModInit(path, importPath)
}

func generateReadmeAndChanges(path, importPath, apiName string) error {
	readmePath := filepath.Join(path, "README.md")
	log.Printf("Creating %q", readmePath)
	readmeFile, err := os.Create(readmePath)
	if err != nil {
		return err
	}
	defer readmeFile.Close()
	t := template.Must(template.New("readme").Parse(readmeTmpl))
	readmeData := struct {
		Name       string
		ImportPath string
	}{
		Name:       apiName,
		ImportPath: importPath,
	}
	if err := t.Execute(readmeFile, readmeData); err != nil {
		return err
	}

	changesPath := filepath.Join(path, "CHANGES.md")
	log.Printf("Creating %q", changesPath)
	changesFile, err := os.Create(changesPath)
	if err != nil {
		return err
	}
	defer changesFile.Close()
	t2 := template.Must(template.New("changes").Parse(changesTmpl))
	changesData := struct {
		Package string
	}{
		Package: pkgName(importPath),
	}
	return t2.Execute(changesFile, changesData)
}

func pkgName(importPath string) string {
	ss := strings.Split(importPath, "/")
	return ss[len(ss)-1]
}
