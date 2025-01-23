### Prerequisites

- Go (example: 1.23) installed on your machine.
- A code editor (e.g., VSCode).
- Also docker installed on your machine

### Installation

1. **Clone the Repository**

   ```bash
   git clone git@github.com:balqisgautama/link-shortener-without-persistent-database.git
   cd url-shortener-api
   ```

2. **Run the Server**

   ```bash
   make run
   ```

   The server will start on `http://localhost:8080`.

## API Endpoints

### 1. Create a Shortened URL

**Endpoint:** `POST /shorten`

**Description:** Create a shortened URL with a specified expiry time.

**Request Body:**

```json
{
  "original": "https://www.google.com/",
  "expiry": 60
}
```

**Response:**

```json
{
  "original": "https://www.google.com/",
  "base_url": "http://localhost:8080",
  "shortened": "ICC5zKG8",
  "click_count": 0,
  "expiry": 60,
  "expiry_duration": 60000000000,
  "expired_at": 1737601770,
  "created_at": 1737601667
}
```

### 2. Retrieve Details of a Shortened URL

**Endpoint:** `GET /shorten/{shortened}`

**Description:** Retrieve details of a specific shortened URL.

**Path Parameter:**

- `shortened` (string): The shortened URL identifier.

**Response:**

```json
{
  "original": "https://www.google.com/",
  "base_url": "http://localhost:8080",
  "shortened": "ICC5zKG8",
  "click_count": 5,
  "expiry": 90,
  "expiry_duration": 90000000000,
  "expired_at": 1737601770,
  "created_at": 1737601667
}
```

### 3. Update the Expiry of a Shortened URL

**Endpoint:** `PUT /shorten`

**Description:** Update the expiry time of a specific shortened URL.

**Request Body:**

```json
{
  "shortened": "ICC5zKG8",
  "expiry": 90
}
```

**Response:**

```json
{
  "original": "https://www.google.com/",
  "base_url": "http://localhost:8080",
  "shortened": "ICC5zKG8",
  "click_count": 5,
  "expiry": 90,
  "expiry_duration": 90000000000,
  "expired_at": 1737601770,
  "created_at": 1737601667
}
```

## Models

### 1. CreateShortenedUrlRequest

**Description:** Request body for creating a shortened URL.

**Properties:**

- `original` (string): The original URL to be shortened.
- `expiry` (integer): The duration in seconds until the shortened URL expires.

**Example:**

```json
{
  "original": "https://www.google.com/",
  "expiry": 60
}
```

### 2. UpdateShortenedUrlRequest

**Description:** Request body for updating the expiry of a shortened URL.

**Properties:**

- `shortened` (string): The shortened URL identifier.
- `expiry` (integer): The new expiry duration in seconds.

**Example:**

```json
{
  "shortened": "ICC5zKG8",
  "expiry": 90
}
```

### 3. ShortenedUrlResponse

**Description:** Response body for shortened URL details.

**Properties:**

- `original` (string): The original URL.
- `base_url` (string): The base URL of the service.
- `shortened` (string): The shortened URL identifier.
- `click_count` (integer): The number of times the shortened URL has been clicked.
- `expiry` (integer): The duration in seconds until the shortened URL expires.
- `expiry_duration` (integer): The expiry duration in nanoseconds.
- `expired_at` (integer): The timestamp when the shortened URL will expire.
- `created_at` (integer): The timestamp when the shortened URL was created.

**Example:**

```json
{
  "original": "https://www.google.com/",
  "base_url": "http://localhost:8080",
  "shortened": "ICC5zKG8",
  "click_count": 5,
  "expiry": 90,
  "expiry_duration": 90000000000,
  "expired_at": 1737601770,
  "created_at": 1737601667
}
```

---
```
Please review the content and let me know if you have any feedback or if there are any changes you would like to make!