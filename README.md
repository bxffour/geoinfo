# Geoinfo API

The Geoinfo API is a RESTful API that provides information about countries. It is heavily inspired by the Restcountries API and is written in Golang. The API offers various endpoints to access country data based on different criteria, making it a comprehensive resource for obtaining details about countries around the world.

## Helm Chart Installation

The Geoinfo API can be deployed using Helm, a Kubernetes package manager. To deploy the Geoinfo API with the desired configurations, you can use the provided Helm chart. Below are the steps to install the Geoinfo API using the Helm chart.

### Prerequisites

- Kubernetes cluster is set up and accessible from your command-line environment.
- Helm is installed on your local machine and configured to work with your Kubernetes cluster.

### Step 1: Clone the Repository

Clone the repository that contains the Helm chart and the example `values.yml` file.

```
git clone https://github.com/bxffour/geoinfo
cd deployments/charts/geoinfo
```

### Step 2: Customize the Configuration

Edit the `values.yml` file to specify your desired configurations for the Geoinfo API deployment. The example `values.yml` file should look like the following:

```yaml
geoinfo:
  name: geoinfo
  port: 6996
  image:
    name: "ghcr.io/bxffour/geoinfo/api"
    tag: "0.9.1-alpha"
  configFile: "config.yaml"
  credsFile: "secret.toml"
  config:
    env: staging
    db:
      max-open-conns: 25
      max-idle-conns: 25
      max-idle-time: 35m
    limiter:
      rps: 4.3
      burst: 8
      enabled: true
  database:
    credentials:
      user: "geoinfo"
      password: "mypassword"
      dbname: "geoinfo"
      host: "postgresql.default.svc.cluster.local"
      port: 5432
    tls:
      enable: true
      autoGen: true
      sslkey: "tls.key"
      sslcert: "tls.crt"
      sslrootcert: "ca.pem"
      sslmode: "require"
      
batch:
  name: "pgdata-dump"
  image:
    name: "ghcr.io/bxffour/geoinfo/bootstrap"
    tag: "0.9.2"
  database:
    user: postgres
    password: valkyrie2
```

In the above configuration, you can customize various aspects of the Geoinfo API deployment, such as the container image, database credentials, TLS settings, and more.

### Step 3: Deploy the Geoinfo API

To deploy the Geoinfo API with the Helm chart, use the following Helm command:

```
helm install geoinfo ./ --values values.yml
```

This command will deploy the Geoinfo API with the specified configurations.

### Step 4: Access the Geoinfo API

After the deployment is successful, you can access the Geoinfo API using the specified port (in this case, `6996`) and the API endpoints described in the documentation.

## Batch Container for Database Setup

The Helm chart includes a batch container that sets up the database for the Geoinfo API. The batch container uses the `ghcr.io/bxffour/geoinfo/bootstrap` image with tag `0.9.2`. It ensures that the necessary database schema and initial data are available for the API to function correctly. This setup is essential for the Geoinfo API to provide accurate and up-to-date country information.

By following the steps outlined above, you can successfully deploy the Geoinfo API and utilize its endpoints to access valuable data about countries across the globe!

## API Endpoints

### 1. Get All Countries

Endpoint: `GET /all`

Description: Returns a list of all countries.

Parameters:
- `page` (optional): Page number for pagination.
- `page_size` (optional): Number of items per page for pagination.

Example Request:
```
GET /all?page=1&page_size=10
```

### 2. Get Country by Name

Endpoint: `GET /name/{name}`

Description: Returns a country based on a given name.

Parameters:
- `name` (required): Search for a country by name.

Example Request:
```
GET /name/ghana
```

### 3. Search by Country Code

Endpoint: `GET /code/{code}`

Description: Search for a country by CCA2, CCN3, CCA3, or CIOC country code.

Parameters:
- `code` (required): Country code (CCA2, CCN3, CCA3, or CIOC).

Example Request:
```
GET /code/GH
```

### 4. Search by Multiple Country Codes

Endpoint: `GET /codes`

Description: Search for countries that match a comma-separated list of codes.

Parameters:
- `codes` (required): Comma-separated list of country codes.

Example Request:
```
GET /codes?codes=gha,per,usa
```

### 5. Search by Currency

Endpoint: `GET /currency/{currency}`

Description: Search for countries by currency name or code.

Parameters:
- `currency` (required): Currency name or code.
- `page` (optional): Page number for pagination.
- `page_size` (optional): Number of items per page for pagination.

Example Request:
```
GET /currency/dollar
```

### 6. Search by Demonym

Endpoint: `GET /demonym/{demonym}`

Description: Search for countries by how citizens are called (demonym).

Parameters:
- `demonym` (required): Demonym or how citizens are called.

Example Request:
```
GET /demonym/ghanaian
```

### 7. Search by Language

Endpoint: `GET /lang/{lang}`

Description: Search for countries by the language spoken.

Parameters:
- `lang` (required): Language.

Example Request:
```
GET /lang/english
```

### 8. Search by Capital

Endpoint: `GET /capital/{capital}`

Description: Search for countries by the capital city.

Parameters:
- `capital` (required): Capital city.

Example Request:
```
GET /capital/accra
```

### 9. Search by Region

Endpoint: `GET /region/{region}`

Description: Search for countries by region.

Parameters:
- `region` (required): Region.
- `page` (optional): Page number for pagination.
- `page_size` (optional): Number of items per page for pagination.

Example Request:
```
GET /region/africa?page=1&page_size=10
```

### 10. Search by Subregion

Endpoint: `GET /subregion/{subregion}`

Description: Search for countries by subregion.

Parameters:
- `subregion` (required): Subregion.
- `page` (optional): Page number for pagination.
- `page_size` (optional): Number of items per page for pagination.

Example Request:
```
GET /subregion/west%20africa?page=1&page_size=10
```

### 11. Search by Translation

Endpoint: `GET /translation/{translation}`

Description: Search for countries by the translation of the country name.

Parameters:
- `translation` (required): Translation of the country name.

Example Request:
```
GET /translation/alemania
```

## Example Response

The API response will contain detailed information about the requested country, including its official and common names, currencies, languages, geographic location, population, and more.

For detailed information on the response schema, please refer to the OpenAPI specification provided in the API documentation.

## API Base URL

The API base URL is `http://localhost:8080/v1/countries`. You can make requests to the above-mentioned endpoints by appending them to the base URL.

Feel free to explore the Geoinfo API and utilize its endpoints to access valuable data about countries across the globe!
