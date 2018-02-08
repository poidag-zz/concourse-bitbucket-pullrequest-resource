package models

import "time"

// Author represents a real person, referenced as either an "Author" or a "User"
type Author struct {
	DisplayName string `json:"display_name"`
	Links       struct {
		Avatar struct {
			Href string `json:"href"`
		} `json:"avatar"`
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"links,omitempty"`
	Type     string `json:"type"`
	Username string `json:"username"`
	UUID     string `json:"uuid"`
}

// Links is the structure of links and references attached to many Bitbucket API responses.
type Links struct {
	Activity struct {
		Href string `json:"href"`
	} `json:"activity"`
	Approve struct {
		Href string `json:"href"`
	} `json:"approve"`
	Avatar struct {
		Href string `json:"href"`
	} `json:"avatar"`
	Comments struct {
		Href string `json:"href"`
	} `json:"comments"`
	Commits struct {
		Href string `json:"href"`
	} `json:"commits"`
	Decline struct {
		Href string `json:"href"`
	} `json:"decline"`
	Diff struct {
		Href string `json:"href"`
	} `json:"diff"`
	HTML struct {
		Href string `json:"href"`
	} `json:"html"`
	Merge struct {
		Href string `json:"href"`
	} `json:"merge"`
	Self struct {
		Href string `json:"href"`
	} `json:"self"`
	Statuses struct {
		Href string `json:"href"`
	} `json:"statuses"`
}

type PagedGenericResponse struct {
	Page     int               `json:"page"`
	Pagelen  int               `json:"pagelen"`
	Size     int               `json:"size"`
	Next     string            `json:"next"`
	Previous string            `json:"previous"`
	Values   []GenericResponse `json:"values"`
}

// GenericResponse is the standard JSON response from the Bitbucket API covering most object types.
type GenericResponse struct {
	Author            Author         `json:"author,omitempty"`
	CloseSourceBranch bool           `json:"close_source_branch,omitempty"`
	ClosedBy          interface{}    `json:"closed_by,omitempty"`
	CommentCount      int            `json:"comment_count,omitempty"`
	Content           CommentContent `json:"content,omitempty"`
	CreatedOn         time.Time      `json:"created_on"`
	Description       string         `json:"description"`
	Destination       struct {
		Branch struct {
			Name string `json:"name,omitempty"`
		} `json:"branch,omitempty"`
		Commit struct {
			Hash  string `json:"hash"`
			Links Links  `json:"links"`
		} `json:"commit,omitempty"`
		Repository struct {
			FullName string `json:"full_name"`
			Links    Links  `json:"links,omitempty"`
			Name     string `json:"name"`
			Type     string `json:"type"`
			UUID     string `json:"uuid,omitempty"`
		} `json:"repository,omitempty"`
	} `json:"destination,omitempty"`
	ID     int `json:"id"`
	Inline *struct {
		Path string `json:"path,omitempty"`
	} `json:"inline,omitempty"`
	Links       Links       `json:"links,omitempty"`
	MergeCommit interface{} `json:"merge_commit,omitempty"`
	Parent      *struct {
		ID int `json:"id,omitempty"`
	} `json:"parent,omitempty"`
	Participants []struct {
		Approved bool   `json:"approved"`
		Role     string `json:"role"`
		Type     string `json:"type"`
		User     struct {
			DisplayName string `json:"display_name"`
			Links       struct {
				Avatar struct {
					Href string `json:"href"`
				} `json:"avatar"`
				HTML struct {
					Href string `json:"href"`
				} `json:"html"`
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"links"`
			Type     string `json:"type"`
			Username string `json:"username"`
			UUID     string `json:"uuid"`
		} `json:"user"`
	} `json:"participants"`
	Pullrequest struct {
		Type  string `json:"type,omitempty"`
		ID    int    `json:"id,omitempty"`
		Links Links  `json:"links,omitempty"`
	} `json:"pullrequest,omitempty"`
	Reason string `json:"reason,omitempty"`
	Source struct {
		Branch struct {
			Name string `json:"name"`
		} `json:"branch,omitempty"`
		Commit struct {
			Hash  string `json:"hash"`
			Links Links  `json:"links"`
		} `json:"commit,omitempty"`
		Repository struct {
			FullName string `json:"full_name"`
			Links    Links  `json:"links,omitempty"`
			Name     string `json:"name"`
			Type     string `json:"type"`
			UUID     string `json:"uuid"`
		} `json:"repository,omitempty"`
	} `json:"source,omitempty"`
	State     string    `json:"state,omitempty"`
	TaskCount int       `json:"task_count,omitempty"`
	Title     string    `json:"title,omitempty"`
	Type      string    `json:"type,omitempty"`
	UpdatedOn time.Time `json:"updated_on"`
	User      Author    `json:"user,omitempty"`
}

// CheckRequest is the struct/JSON that is supplied to "check", coming from the Concourse pipeline under "resources"
type CheckRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

// CheckResponse is the struct/JSON returned from "check"
type CheckResponse []Version

// Params ... (referenced from OutRequest)
type Params struct {
	State       string `json:"state"`
	PullRequest string `json:"pull_request"`
	Commit      string `json:"commit"`
}

// Source ... (referenced from CheckRequest)
type Source struct {
	Repo         string `json:"repo"`
	Secret       string `json:"secret"`
	Key          string `json:"key"`
	Team         string `json:"team"`
	URL          string `json:"url"`
	APIVersion   string `json:"version"`
	ConcourseURL string `json:"concourse_url"`
}

// Version ... (referenced from CheckRequest)
type Version struct {
	Commit      string `json:"commit"`
	PullRequest string `json:"pullrequest"`
	Link        string `json:"link,omitempty"`
}

// InRequest is the struct/JSON supplied as input to "in" - Concourse pipeline "get"
type InRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

// InResponse is the struct/JSON that is output from "in".
type InResponse struct {
	Version  Version  `json:"version"`
	Metadata Metadata `json:"metadata"`
}

// Metadata holds multiple MetadataField, which is output from "in" and "out"
type Metadata []MetadataField

// MetadataField holding data presented as metedata in "in" and "out"
type MetadataField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// OutRequest is the struct/JSON supplied as input to "out" - Concourse pipeline "put"
type OutRequest struct {
	Params  Params  `json:"params"`
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

// OutResponse is the struct/JSON that is output from "out".
type OutResponse struct {
	Version  Version  `json:"version"`
	Metadata Metadata `json:"metadata"`
}

// OutStatus holds data about a build's status.
type OutStatus struct {
	State string `json:"state"`
	Key   string `json:"key"`
	URL   string `json:"url"`
}

// type CredentialsRequest2 struct {
// 	GrantType string `json:"grant_type"`
// }

// Token holds Authentication Tokens for accessing the Bitbucket API.
type Token struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scopes       string `json:"scopes"`
	TokenType    string `json:"token_type"`
}

// CommitResponse represents the Commit Status response from the Bitbucket API.
// <https://developer.atlassian.com/bitbucket/api/2/reference/resource/repositories/%7Busername%7D/%7Brepo_slug%7D/commit/%7Bnode%7D/statuses>
type CommitResponse struct {
	Page    int `json:"page"`
	Pagelen int `json:"pagelen"`
	Size    int `json:"size"`
	Values  []struct {
		CreatedOn   time.Time   `json:"created_on"`
		Description string      `json:"description"`
		Key         string      `json:"key"`
		Links       Links       `json:"links"`
		Name        string      `json:"name"`
		Refname     interface{} `json:"refname"`
		Repository  struct {
			FullName string `json:"full_name"`
			Links    struct {
				Avatar struct {
					Href string `json:"href"`
				} `json:"avatar"`
				HTML struct {
					Href string `json:"href"`
				} `json:"html"`
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"links"`
			Name string `json:"name"`
			Type string `json:"type"`
			UUID string `json:"uuid"`
		} `json:"repository"`
		State     string `json:"state"`
		Type      string `json:"type"`
		UpdatedOn string `json:"updated_on"`
		URL       string `json:"url"`
	} `json:"values"`
}

// CommentContent is the actual text of a comment.
type CommentContent struct {
	Raw  string `json:"raw"`
	HTML string `json:"html"`
}

// Comment represents a comment on a Pull Request.
type Comment struct {
	// Pullrequest struct{}       `json:"pullrequest"`
	Content   CommentContent `json:"content"`
	CreatedOn time.Time      `json:"created_on"`

	User   Author `json:"user"`
	Inline struct {
		Path string `json:"path"`
	} `json:"inline"`
	UpdatedOn time.Time `json:"updated_on"`
	Type      string    `json:"type"`
	ID        int       `json:"id"`
	Link      string
}
