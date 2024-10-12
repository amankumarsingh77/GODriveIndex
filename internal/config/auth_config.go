package config

type AuthConfig struct {
	SiteName                   string `json:"siteName"`
	ClientID                   string `json:"client_id"`
	ClientSecret               string `json:"client_secret"`
	RefreshToken               string `json:"refresh_token"`
	ServiceAccount             bool   `json:"service_account"`
	ServiceAccountJSON         string `json:"service_account_json"`
	FilesListPageSize          int    `json:"files_list_page_size"`
	SearchResultListPageSize   int    `json:"search_result_list_page_size"`
	EnableCorsFileDown         bool   `json:"enable_cors_file_down"`
	EnablePasswordFileVerify   bool   `json:"enable_password_file_verify"`
	DirectLinkProtection       bool   `json:"direct_link_protection"`
	DisableAnonymousDownload   bool   `json:"disable_anonymous_download"`
	FileLinkExpiry             int    `json:"file_link_expiry"`
	SearchAllDrives            bool   `json:"search_all_drives"`
	EnableLogin                bool   `json:"enable_login"`
	EnableSignup               bool   `json:"enable_signup"`
	EnableSocialLogin          bool   `json:"enable_social_login"`
	GoogleClientIDForLogin     string `json:"google_client_id_for_login"`
	GoogleClientSecretForLogin string `json:"google_client_secret_for_login"`
	RedirectDomain             string `json:"redirect_domain"`
	LoginDatabase              string `json:"login_database"`
	LoginDays                  int    `json:"login_days"`
	EnableIPLock               bool   `json:"enable_ip_lock"`
	SingleSession              bool   `json:"single_session"`
	IPChangedAction            bool   `json:"ip_changed_action"`
	UsersList                  []User `json:"users_list"`
	Roots                      []Root `json:"roots"`
}

type ServiceAccount struct {
	Type         string `json:"type"`
	ProjectID    string `json:"project_id"`
	PrivateKeyID string `json:"private_key_id"`
	PrivateKey   string `json:"private_key"`
	ClientEmail  string `json:"client_email"`
	ClientID     string `json:"client_id"`
	AuthURI      string `json:"auth_uri"`
	TokenURI     string `json:"token_uri"`
	AuthProvider string `json:"auth_provider_x509_cert_url"`
	ClientUrl    string `json:"client_x509_cert_url"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Root struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	ProtectFileLink bool   `json:"protect_file_link"`
}
