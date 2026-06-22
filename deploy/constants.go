package deploy

const (
	HealthPath = "/health"
)

const testComposeTemplate = `services:
${APP_SERVICE_BLOCK}
${LOCAL_POSTGRES_SERVICE_BLOCK}

networks:
  quollix_${APP_NAME}:
    name: quollix_${APP_NAME}
`

const standalonePostgresComposeTemplate = `services:
${LOCAL_POSTGRES_SERVICE_BLOCK}

networks:
  quollix_${APP_NAME}:
    name: quollix_${APP_NAME}
`
