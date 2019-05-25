# Configuring

On Windows:

    copy .env.example .env

Fill the .env file with development settings

# Running

Default:

    go run main.go serve

On Windows (without getting prompted by firewall so many times):

    go build && .\go-api-base.exe serve
