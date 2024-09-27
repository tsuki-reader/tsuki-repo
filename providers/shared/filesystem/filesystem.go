package main

import "github.com/tsuki-reader/nisshoku/providers"

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
	panic("TODO")
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
