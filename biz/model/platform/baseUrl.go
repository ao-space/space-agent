package platform

type PlatformUrl struct {
	AppSSPlatformUrl      string
	AppPSPlatformUrl      string
	AppAppstoreAppApiUrl  string
	AppAppstoreAppSignUrl string
}

var BaseUrlMap map[string]PlatformUrl

func init() {
	baseUrlMap := make(map[string]PlatformUrl)
	baseUrlMap["CN"] = PlatformUrl{
		AppSSPlatformUrl:      "https://ao.space",
		AppPSPlatformUrl:      "https://ao.space",
		AppAppstoreAppApiUrl:  "https://api.apps.ao.space",
		AppAppstoreAppSignUrl: "https://auth.apps.ao.space",
	}
	baseUrlMap["SG"] = PlatformUrl{
		AppSSPlatformUrl:      "https://aospace.sg",
		AppPSPlatformUrl:      "https://aospace.sg",
		AppAppstoreAppApiUrl:  "https://api.apps.aospace.sg",
		AppAppstoreAppSignUrl: "https://auth.apps.aospace.sg",
	}
	baseUrlMap["dev"] = PlatformUrl{
		AppSSPlatformUrl:      "https://dev.eulix.xyz",
		AppPSPlatformUrl:      "https://dev.eulix.xyz",
		AppAppstoreAppApiUrl:  "https://api.dev-apps.eulix.xyz",
		AppAppstoreAppSignUrl: "https://auth.dev-apps.eulix.xyz",
	}
	baseUrlMap["dev2"] = PlatformUrl{
		AppSSPlatformUrl:      "https://dev2.eulix.xyz",
		AppPSPlatformUrl:      "https://dev2.eulix.xyz",
		AppAppstoreAppApiUrl:  "https://api.dev2-apps.eulix.xyz",
		AppAppstoreAppSignUrl: "https://auth.dev2-apps.eulix.xyz",
	}
	baseUrlMap["sit"] = PlatformUrl{
		AppSSPlatformUrl:      "https://sit.eulix.xyz",
		AppPSPlatformUrl:      "https://sit.eulix.xyz",
		AppAppstoreAppApiUrl:  "https://api.sit-apps.eulix.xyz",
		AppAppstoreAppSignUrl: "https://auth.sit-apps.eulix.xyz",
	}
	baseUrlMap["sit2"] = PlatformUrl{
		AppSSPlatformUrl:      "https://sit2.eulix.xyz",
		AppPSPlatformUrl:      "https://sit2.eulix.xyz",
		AppAppstoreAppApiUrl:  "https://api.sit2-apps.eulix.xyz",
		AppAppstoreAppSignUrl: "https://auth.sit2-apps.eulix.xyz",
	}

	BaseUrlMap = baseUrlMap
}
