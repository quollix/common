package validation

const (
	LabelField            = "label"
	ServiceField          = "service"
	VolumeField           = "volume"
	VolumeNameField       = "volume_name"
	MaintainerField       = "maintainer"
	AppField              = "app"
	KeyField              = "key"
	ImageField            = "image"
	ActualContainerName   = "actual_container_name"
	ExpectedContainerName = "expected_container_name"
	PortField             = "port"
	ExpectedPrefixField   = "expected_prefix"
	ExpectedPartCount     = "expected_part_count"
	ActualPartCount       = "actual_part_count"
)

var (
	notAllowedTopLevelKeyword             = "not allowed root keyword in docker-compose.yml"
	notAllowedKeyInService                = "not allowed key in service"
	mainServiceMustBeDefined              = "the main service was not defined"
	serviceNeedsContainerNameKeyword      = "service must set 'container_name' keyword"
	wrongContainerNameValue               = "service has invalid container_name, it should be '<maintainer>_<app>_<service>'"
	unexpectedUnderscoreCountContainer    = "unexpected number of underscores in container_name"
	mustSetTheDockerImageTag              = "the image tag must be set, like 'gitea/gitea:10.5'"
	notAllowedLatestDockerImageTag        = "the 'latest' tag is forbidden, to get reproducible apps, only fixed tags with specific software version should be used"
	SystemAppNamesAreAlreadyReserved      = "system app names are not allowed"
	devicesKeywordIsForbidden             = "'devices' keyword is not allowed"
	deployKeywordMustOnlyContainResources = "'deploy' keyword must only contain 'resources' keyword"
	globalVolumeShouldNotHaveSubKeywords  = "volume has sub-keywords, which are not allowed"
	mustSetImageKey                       = "service must have 'image' key"

	exposingDefaultPortIsForbidden = "exposing this port is forbidden"

	hostDirectoriesMountedForbidden = "host directories are mounted which is forbidden"

	volumeNamePrefixIsWrong          = "volume name prefix is wrong"
	volumeEntryMissingColonSeparator = "volume entry is missing ':' separator"
	unexpectedUnderscoreCountVolume  = "unexpected number of underscores in volume name"
)
