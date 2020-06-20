package tagging

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"net/http"
	"testing"
)

func TestValidateTagNotExisting_Github(t *testing.T) {
	// Testing tag not existing
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.github.com").
		Get("/repos/repo/git/refs/tags").
		Reply(http.StatusNotFound)
	assertTest := assert.New(t)
	repo := GithubProperties{RepoProperties{"username", "password", "repo", "tag", "hash", "", ""}}
	results := repo.ValidateTag()
	assertTest.True(results.TagDoesntExist)
	assertTest.False(results.TagExistsWithProvidedHash)
}

func TestValidateTagUnauthorized_Githhub(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.github.com").
		Get("/repos/repo/git/refs/tags").
		Reply(http.StatusUnauthorized)
	assertTest := assert.New(t)
	// Testing a 403
	repo := GithubProperties{RepoProperties{"username", "password", "repo", "tag", "hash", "", ""}}
	results := repo.ValidateTag()
	assertTest.False(results.TagDoesntExist)
	assertTest.False(results.TagExistsWithProvidedHash)
}

func TestValidateTagExistingSameHash_Github(t *testing.T) {
	target := Object{Sha: "hash"}
	tag := Tag{Object: target}
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.github.com").
		Get("/repos/repo/git/refs/tags").
		Reply(http.StatusOK).
		JSON(tag)

	assertTest := assert.New(t)
	// Testing 200 response and hash is the same
	repo := GithubProperties{RepoProperties{"username", "password", "repo", "tag", "hash", "", ""}}
	results := repo.ValidateTag()
	assertTest.False(results.TagDoesntExist)
	assertTest.True(results.TagExistsWithProvidedHash)
}

func TestValidateTagExistingMismatchHash_Github(t *testing.T) {
	assertTest := assert.New(t)
	// Testing 200 response but hash is not the same
	target := Object{Sha: "hash"}
	tag := Tag{Object: target}
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.github.com").
		Get("/repos/repo/git/refs/tags").
		Reply(http.StatusOK).
		JSON(tag)
	repo := GithubProperties{RepoProperties{"username", "password", "repo", "tag", "not_hash", "", ""}}
	results := repo.ValidateTag()
	assertTest.False(results.TagDoesntExist)
	assertTest.False(results.TagExistsWithProvidedHash)
}

func TestValidateTagOtherError_Github(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.github.com").
		Get("/repos/repo/git/refs/tags").
		Reply(http.StatusServiceUnavailable)
	assertTest := assert.New(t)
	// Testing a 403
	repo := GithubProperties{RepoProperties{"username", "password", "repo", "tag", "hash", "", ""}}
	results := repo.ValidateTag()
	assertTest.False(results.TagDoesntExist)
	assertTest.False(results.TagExistsWithProvidedHash)
}

func TestCreateTagNotFound_Github(t *testing.T) {
	// Testing tag not existing
	target := Object{Sha: "tag"}
	tag := Tag{Object: target}
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.github.com").
		Get("/repos/repo/git/refs/tags").
		Reply(http.StatusNotFound)

	gock.New("https://api.github.com").
		Post("/repos/repo/releases").
		Reply(http.StatusNotFound).
		JSON(tag)

	assertTest := assert.New(t)
	repo := GithubProperties{RepoProperties{"username", "password", "repo", "tag", "hash", "", ""}}
	assertTest.False(repo.CreateTag())
}

func TestCreateTagUnauthorized_Github(t *testing.T) {
	// Testing a 401
	body := GithubRelease{TargetCommitish: "hash", Prerelease: false, Draft: false, Body: "hello", TagName: "tag", Name: "tag"}
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.github.com").
		Get("/repos/repo/git/refs/tags").
		Reply(http.StatusNotFound)

	gock.New("https://api.github.com").
		Post("/repos/repo/releases").
		Reply(http.StatusUnauthorized).
		JSON(body)
	assertTest := assert.New(t)
	repo := GithubProperties{RepoProperties{"username", "password", "repo", "tag", "hash", "", "hello"}}
	assertTest.False(repo.CreateTag())
}

func TestCreateTagSuccessful_Github(t *testing.T) {
	// Testing 201 response
	body := GithubRelease{TargetCommitish: "hash", Prerelease: false, Draft: false, Body: "hello", TagName: "tag", Name: "tag"}
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.github.com").
		Get("/repos/repo/git/refs/tags").
		Reply(http.StatusNotFound)

	gock.New("https://api.github.com").
		Post("/repos/repo/releases").
		Reply(http.StatusCreated).
		JSON(body)
	assertTest := assert.New(t)
	repo := GithubProperties{RepoProperties{"username", "password", "repo", "tag", "hash", "", "hello"}}
	assertTest.True(repo.CreateTag())
}

func TestCreateTagSuccessfulWithHostOverride_Github(t *testing.T) {
	// Testing 201 response
	body := GithubRelease{TargetCommitish: "hash", Prerelease: false, Draft: false, Body: "hello", TagName: "tag", Name: "tag"}
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.personal-github.com").
		Get("/repos/repo/git/refs/tags").
		Reply(http.StatusNotFound)

	gock.New("https://api.personal-github.com").
		Post("/repos/repo/releases").
		Reply(http.StatusCreated).
		JSON(body)
	assertTest := assert.New(t)
	repo := GithubProperties{RepoProperties{"username", "password", "repo", "tag", "hash", "https://api.personal-github.com", "hello"}}
	assertTest.True(repo.CreateTag())
}

func TestCreateTagAlreadyExists_Github(t *testing.T) {
	target := Object{Sha: "hash"}
	tag := Tag{Object: target}
	defer gock.Off() // Flush pending mocks after test execution
	gock.New("https://api.github.com").
		Get("/repos/repo/git/refs/tags").
		Reply(http.StatusOK).
		JSON(tag)
	assertTest := assert.New(t)
	repo := GithubProperties{RepoProperties{"username", "password", "repo", "test", "hash", "", "hello"}}
	assertTest.True(repo.CreateTag())
}

func TestCreateError_Github(t *testing.T) {
	target := Object{Sha: "hash"}
	tag := Tag{Object: target}
	errorMessage := GithubError{Code: "already_exists"}
	response := GithubBadResponse{Errors: []GithubError{errorMessage}}
	defer gock.Off() // Flush pending mocks after test execution
	gock.New("https://api.github.com").
		Get("/repos/repo/git/refs/tags").
		Reply(http.StatusNotFound).
		JSON(tag)
	gock.New("https://api.github.com").
		Post("/repos/repo/releases").
		Reply(http.StatusUnprocessableEntity).
		JSON(response)
	assertTest := assert.New(t)
	repo := GithubProperties{RepoProperties{"username", "password", "repo", "test", "hash", "", "hello"}}
	assertTest.False(repo.CreateTag())
}

func TestCreateTagOtherError_Github(t *testing.T) {
	errorMessage := GithubError{Code: "blah"}
	response := GithubBadResponse{Errors: []GithubError{errorMessage}}
	defer gock.Off() // Flush pending mocks after test execution
	gock.New("https://api.github.com").
		Get("/repos/repo/git/refs/tags").
		Reply(http.StatusOK)
	gock.New("https://api.github.com").
		Post("/repos/repo/releases").
		Reply(http.StatusBadRequest).
		JSON(response)
	assertTest := assert.New(t)
	// Testing 400 response has been created, should never happen if validate is called first
	repo := GithubProperties{RepoProperties{"username", "password", "repo", "tag", "hash", "", "hello"}}
	assertTest.False(repo.CreateTag())
}

func TestCreateTagOtherErrorResponse_Github(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution
	gock.New("https://api.github.com").
		Get("/repos/repo/git/refs/tags").
		Reply(http.StatusOK)
	gock.New("https://api.github.com").
		Post("/repos/repo/releases").
		Reply(http.StatusServiceUnavailable)
	assertTest := assert.New(t)
	// Testing 400 response has been created, should never happen if validate is called first
	repo := GithubProperties{RepoProperties{"username", "password", "repo", "tag", "hash", "", "hello"}}
	assertTest.False(repo.CreateTag())
}
