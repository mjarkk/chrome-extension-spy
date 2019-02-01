package firefox

// FF is the main package
type FF struct {
	LaunchCMD           string // this is something like "firefox" on linux and something like "C:\Program Files\firefox\firefox.exe" on windows
	UserProfileLocation string // the default user profile location
	TmpDirs             FfTmpDirs
	HasErr              error // some type function had a error and will be stored here
}

// FfTmpDirs is the temp dirs struct for the FF type
type FfTmpDirs struct {
	UnpackExts string // The unpacked extensions will be here
	Profile    string // The created user profile
}

// Err return true if HasErr is not nil
func (f *FF) Err() bool {
	return f.HasErr != nil
}

// Setup returns the default firefox struct
func Setup() FF {
	f := FF{
		LaunchCMD:           "",
		UserProfileLocation: "",
		TmpDirs: FfTmpDirs{
			UnpackExts: "",
			Profile:    "",
		},
		HasErr: nil,
	}

	f.GetLaunchCMD()
	f.GetUserLocation()
	f.GetRawExts()
	f.CreateEmptyProfile()
	f.PackExtensions()

	return f
}
