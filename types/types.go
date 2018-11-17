package types

// SmallRequest is a small version of request
type SmallRequest struct {
	Pkg  string `json:"pkg"`  // package
	URL  string `json:"url"`  // the request url
	Type string `json:"type"` // request type (POST or GET)
	Code int    `json:"code"` // status code
	Hash string `json:"hash"` // the request hash
}

// Request is the content of a chrome request
type Request struct {
	Extension       ChromeExtension   `json:"extension"`       // details of the chrome extension
	URL             string            `json:"url"`             // request url
	Type            string            `json:"type"`            // request type (POST or GET)
	StatusCode      int               `json:"statusCode"`      // response status code
	RequestHeaders  map[string]string `json:"requestHeaders"`  // request headers
	ResponseHeaders map[string]string `json:"responseHeaders"` // response headers
	PostBody        string            `json:"postBody"`        // post request body
	ResData         string            `json:"resData"`         // response data
	ResData64       string            `json:"resData64"`       // raw response in base64
	Hash            string            `json:"hash"`            // the request hash
}

// FullAndSmallExt has both the ChromeExtension and ExtensionManifest in 1 struct
type FullAndSmallExt struct {
	Small ChromeExtension
	Full  ExtensionManifest
}

// ChromeExtension is a small version of extensionManifest with just the right amound of data
type ChromeExtension struct {
	Pkg         string `json:"pkg"`         // the chrome extension folder name
	PkgVersion  string `json:"pkgVersion"`  // the package version
	FullPkgURL  string `json:"fullPkgURL"`  // full url to extension
	Name        string `json:"name"`        // extension name
	ShortName   string `json:"shortName"`   // extension longname
	HomepageURL string `json:"homepageURL"` // homepage url
}

// ExtensionManifest fully covers most manifest.json files from extensions
type ExtensionManifest struct {
	Applications struct {
		Gecko struct {
			ID               string `json:"id"`
			StrictMinVersion string `json:"strict_min_version"`
		} `json:"gecko"`
	} `json:"applications"`
	App struct {
		Launch struct {
			Container string `json:"container"`
			WebURL    string `json:"web_url"`
		} `json:"launch"`
		WebContent struct {
			Enabled bool   `json:"enabled"`
			Origin  string `json:"origin"`
		} `json:"web_content"`
		Urls []string `json:"urls"`
	} `json:"app"`
	Commands struct {
		ExecutePageAction struct {
			SuggestedKey struct {
				Default string `json:"default"`
			} `json:"suggested_key"`
		} `json:"_execute_page_action"`
		DevtoolsBottom struct {
			Description string `json:"description"`
		} `json:"devtools-bottom"`
		DevtoolsLeft struct {
			Description string `json:"description"`
		} `json:"devtools-left"`
		DevtoolsRemote struct {
			Description string `json:"description"`
		} `json:"devtools-remote"`
		DevtoolsRight struct {
			Description string `json:"description"`
		} `json:"devtools-right"`
	} `json:"commands"`
	Background struct {
		Persistent    bool     `json:"persistent"`
		AllowJsAccess bool     `json:"allow_js_access"`
		Scripts       []string `json:"scripts"`
	} `json:"background"`
	ChromeURLOverrides struct {
		Newtab string `json:"newtab"`
	} `json:"chrome_url_overrides"`
	DisplayInLauncher   bool `json:"display_in_launcher"`
	DisplayInNewTabPage bool `json:"display_in_new_tab_page"`
	BrowserAction       struct {
		DefaultIcon  string `json:"default_icon"`
		DefaultTitle string `json:"default_title"`
		DefaultPopup string `json:"default_popup"`
	} `json:"browser_action"`
	ContentScripts []struct {
		AllFrames       bool     `json:"all_frames"`
		CSS             []string `json:"css"`
		Js              []string `json:"js"`
		MatchAboutBlank bool     `json:"match_about_blank"`
		ExcludeGlobs    []string `json:"exclude_globs"`
		Matches         []string `json:"matches"`
		RunAt           string   `json:"run_at"`
	} `json:"content_scripts"`
	ContentSecurityPolicy string `json:"content_security_policy"`
	Description           string `json:"description"`
	Icons                 struct {
		Num16  string `json:"16"`
		Num32  string `json:"32"`
		Num48  string `json:"48"`
		Num64  string `json:"64"`
		Num128 string `json:"128"`
	} `json:"icons"`
	Key                   string      `json:"key"`
	Author                interface{} `json:"author"`
	Incognito             string      `json:"incognito"`
	DevtoolsPage          string      `json:"devtools_page"`
	HomepageURL           string      `json:"homepage_url"`
	ManifestVersion       int         `json:"manifest_version"`
	DefaultLocale         string      `json:"default_locale"`
	OfflineEnabled        bool        `json:"offline_enabled"`
	MinimumChromeVersion  string      `json:"minimum_chrome_version"`
	MinimumOperaVersion   string      `json:"minimum_opera_version"`
	ShortName             string      `json:"short_name"`
	Name                  string      `json:"name"`
	ExternallyConnectable struct {
		Ids     []string `json:"ids"`
		Matches []string `json:"matches"`
	} `json:"externally_connectable"`
	Storage struct {
		ManagedSchema string `json:"managed_schema"`
	} `json:"storage"`
	OptionsPage string        `json:"options_page"`
	Permissions []interface{} `json:"permissions"`
	UpdateURL   string        `json:"update_url"`
	Version     string        `json:"version"`
	Oauth2      struct {
		AutoApprove bool     `json:"auto_approve"`
		ClientID    string   `json:"client_id"`
		Scopes      []string `json:"scopes"`
	} `json:"oauth2"`
	OptionsUI struct {
		ChromeStyle bool   `json:"chrome_style"`
		OpenInTab   bool   `json:"open_in_tab"`
		Page        string `json:"page"`
	} `json:"options_ui"`
	Sandbox struct {
		ContentSecurityPolicy string   `json:"content_security_policy"`
		Pages                 []string `json:"pages"`
	} `json:"sandbox"`
	URLHandlers struct {
		PostmanCollection struct {
			Matches []string `json:"matches"`
			Title   string   `json:"title"`
		} `json:"postman_collection"`
	} `json:"url_handlers"`
	PageAction struct {
		DefaultIcon  string `json:"default_icon"`
		DefaultPopup string `json:"default_popup"`
		DefaultTitle string `json:"default_title"`
	} `json:"page_action"`
	WebAccessibleResources []string `json:"web_accessible_resources"`
}
