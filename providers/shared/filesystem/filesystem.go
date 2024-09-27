package main

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/tsuki-reader/nisshoku/providers"
)

type FilesystemProvider struct {
	libraryPath  string
	providerType providers.ProviderType
}

// NewProvider must take an instance of ProviderParams and return a Provider
func NewProvider(params providers.ProviderParams) providers.Provider {
	var provider FilesystemProvider

	provider.providerType = params.ProviderType
	switch provider.providerType {
	case providers.Comic:
		provider.libraryPath = params.ComicLibraryPath
	case providers.Manga:
		provider.libraryPath = params.MangaLibraryPath
	}

	return &provider
}

func (p *FilesystemProvider) Search(query string) ([]providers.ProviderResult, error) {
	var results []providers.ProviderResult

	filepath.WalkDir(p.libraryPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			name := d.Name()
			if strings.Contains(strings.ToLower(name), strings.ToLower(query)) {
				var result = providers.ProviderResult{
					Title:    name,
					ID:       name,
					Provider: "filesystem",
				}
				results = append(results, result)
			}
		}

		return nil
	})

	return results, nil
}

func (p *FilesystemProvider) GetChapters(id string) ([]providers.Chapter, error) {
	panic("TODO")
}

func (p *FilesystemProvider) GetChapterPages(id string) ([]providers.Page, error) {
	panic("TODO")
}

func (p *FilesystemProvider) ImageHeaders() map[string]string {
	return map[string]string{}
}

func (p *FilesystemProvider) ProviderType() providers.ProviderType {
	return p.providerType
}
