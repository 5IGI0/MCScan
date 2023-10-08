# MCScan API

Welcome to the MCScan API repository, a Go-based API that allows you to search for Minecraft servers based on various criteria, including server version, installed mods, and server description. This API provides valuable information about Minecraft servers, making it easier for you to find the servers that match your preferences.

_Please note that while the name suggests scanning, this API does perform server scans, but it does so primarily for the purpose of updating server information. Servers still need to be added manually to the API._

## Endpoints

### Search Minecraft Servers

Use the following endpoint to search for Minecraft servers based on specific parameters:

**GET**  ``/v1/search``

Available parameters for searching include:

- `version`: Protocol version (Refer to [this](https://github.com/PrismarineJS/minecraft-data/blob/master/data/pc/common/protocolVersions.json))
- `mods`: Mods' slugs, separated by commas.
- `text`: Any text found in the server's MOTD (Message of the Day).

### Get Server Favicon

Retrieve a server's favicon using its `favicon_id`, which can be obtained from the server's object.

**GET** ``/v1/favicon/{favicon_id}``

### Get Server Information

Retrieve detailed information about a server using either its `id` or `addr`. \
_Please note that to obtain aliases in the returned object, you should use the `full=1` GET parameter._

**GET** ``/v1/server/{id|addr}``

## Environment Variables

To configure the MCScan API, you can set the following environment variables:

- `MYSQL_DB_URL`: The URI for the data source, which can be configured following the guidelines provided in the [source data](https://github.com/go-sql-driver/mysql/#dsn-data-source-name) documentation.
- `LISTEN_ADDR`: The address where the HTTP server should listen. You can use unix: to specify a Unix socket.
- `ADD_SERVER_PASSWORD`: A password required to access the /v1/add_servers endpoint for adding servers.

If you have any questions or need assistance, please feel free to reach out.