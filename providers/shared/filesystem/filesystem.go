package main

import (
	"strings"

	"github.com/tsuki-reader/nisshoku/providers"
)

type FilesystemProvider struct {
	libraryPath  string
	providerType providers.ProviderType
	context      providers.ProviderContext
}

// NewProvider must take an instance of ProviderParams and return a Provider
func NewProvider(context providers.ProviderContext) providers.Provider {
	var provider FilesystemProvider

	provider.context = context
	provider.providerType = context.ProviderType
	switch provider.providerType {
	case providers.Comic:
		provider.libraryPath = context.ComicLibraryPath
	case providers.Manga:
		provider.libraryPath = context.MangaLibraryPath
	}

	return &provider
}

func (p *FilesystemProvider) Search(query string) ([]providers.ProviderResult, error) {
	var results []providers.ProviderResult

	files, err := p.context.WalkLibrary(p.libraryPath)
	if err != nil {
		return []providers.ProviderResult{}, err
	}

	for _, f := range files {
		if f.IsDir && strings.Contains(strings.ToLower(f.Name), strings.ToLower(query)) {
			result := providers.ProviderResult{
				Title:            f.Name,
				ID:               f.Fullpath,
				Provider:         "filesystem",
				AlternativeNames: []string{},
				StartYear:        0,
			}

			results = append(results, result)
		}
	}

	return results, nil
}

func (p *FilesystemProvider) GetChapters(id string) ([]providers.Chapter, error) {
	var results []providers.Chapter

	chapterDirs, err := p.context.WalkLibrary(id)
	if err != nil {
		return []providers.Chapter{}, err
	}

	for i, chapterDir := range chapterDirs {
		if chapterDir.IsDir {
			chapter := providers.Chapter{
				Title:          chapterDir.Name,
				ID:             chapterDir.Fullpath,
				Provider:       "filesystem",
				AbsoluteNumber: i,
			}
			results = append(results, chapter)
		}
	}

	return results, nil
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
