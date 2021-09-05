# API documentation

Termora has a basic REST API for querying terms and explanations.
The root endpoint of the API is `https://api.termora.org/v1/`.

API endpoints, response fields, and query parameters may be added, but not removed, without a major version change.

Requests are rate limited to 5 requests per host per second.

## Models

The following three models (usually represented in JSON format) represent the objects in Termora's API.
A `?` after the column type indicates an optional (nullable) parameter.

### Term object

| Key              | Type      | Notes                                                           |
| ---------------- | --------- | --------------------------------------------------------------- |
| id               | number    | The internal numeric ID.                                        |
| category_id      | number    | The category's internal ID.                                     |
| category         | string    |                                                                 |
| name             | string    |                                                                 |
| aliases          | string[]  | Referred to as "synonyms" in the bot.                           |
| description      | string    |                                                                 |
| note             | string?   |                                                                 |
| source           | string    |                                                                 |
| created          | datetime  |                                                                 |
| last_modified    | datetime  | Will be the same as `created` if the term hasn't been modified. |
| tags             | string[]? |                                                                 |
| content_warnings | string?   |                                                                 |
| image_url        | string?   | Currently unused.                                               |
| flags            | number    | A bitmask of term flags.                                        |
| rank             | number?   | Only returned in searches.                                      |
| headline         | string?   | Only returned in searches.                                      |

#### Term flags

| Flag     | Meaning                                            |
| -------- | -------------------------------------------------- |
| `1 << 0` | Hidden from search                                 |
| `1 << 1` | Not shown in random results                        |
| `1 << 2` | Shows a warning if looked up on the bot or website |
| `1 << 3` | Hidden from lists (including the website)          |
| `1 << 4` | Shows a "disputed" note                            |

### Category object

| Key  | Type   | Notes                       |
| ---- | ------ | --------------------------- |
| id   | number | The category's internal ID. |
| name | string |                             |

### Explanation object

| Key         | Type     | Notes                          |
| ----------- | -------- | ------------------------------ |
| id          | number   | The explanation's internal ID. |
| name        | string   |                                |
| aliases     | string[] | Alternative names/triggers.    |
| description | string   |                                |
| created     | datetime |                                |

## Endpoints

### `GET /term/:id`

Gets a term by its numeric ID. Returns a [term object](#term-object) on success,
`404 Not Found` if the term wasn't found,
and `400 Bad Request` if `:id` was not an integer.

**Example request**

```
GET https://api.termora.org/v1/term/1
```

**Example response**

```json
{
    "id": 1,
    "category_id": 1,
    "category": "Plurality",
    "name": "Plural",
    "aliases": [
        "Plurality"
    ],
    "description": "An umbrella term encompassing all phenomena in which multiple consciousnesses cohabit a single brain and body.",
    "source": "https://tulpa.io/terminologies",
    "created": "2020-12-30T15:32:55.354524Z",
    "last_modified": "2021-04-07T13:52:01.993722Z",
    "tags": [
        "Plurality"
    ],
    "flags": 0
}
```

### `GET /search/:term`

Searches the database for a query. Returns an array of [term objects](#term-object) on success,
or `204 No Content` if no results were found.

**Example query**

```
GET https://api.termora.org/v1/search/ace
```

**Example response**

```json
[
    {
        "id": 232,
        "category_id": 2,
        "category": "LGBTQ+",
        "name": "Asexual Spectrum",
        "aliases": [
            "Acespec",
            "Ace-spec"
        ],
        "description": "An umbrella term for any terms that fall under the umbrella term of asexual.",
        "source": "Unknown; already in circulation.",
        "created": "2021-03-19T19:47:40.859011Z",
        "last_modified": "2021-03-26T20:16:41.079386Z",
        "tags": [
            "LGBTQ+",
            "Sexuality"
        ],
        "flags": 0,
        "rank": 0.0833333358168602,
        "headline": "An umbrella term for any terms that fall under the umbrella term of asexual."
    },
    {
        "id": 106,
        "category_id": 2,
        "category": "LGBTQ+",
        "name": "Asexuality",
        "aliases": [
            "Ace"
        ],
        "description": "The lack of sexual attraction. Might also include not being interested in sex, not experiencing a sex drive/libido, or being repulsed by sex.",
        "source": "https://lgbta.wikia.org/wiki/Asexual",
        "created": "2021-02-07T19:11:17.296278Z",
        "last_modified": "2021-03-26T20:16:25.389939Z",
        "tags": [
            "LGBTQ+",
            "Sexuality"
        ],
        "flags": 0,
        "rank": 0.0625,
        "headline": "The lack of sexual attraction. Might also include not being interested in sex, not experiencing"
    },
    // ...
]
```

### `GET /explanations`

Gets all explanations from the database. Returns an array of [explanation objects](#explanation-object) on success,
or `204 No Content` if there are no explanations.

**Example query**

```
GET https://api.termora.org/v1/explanations
```

**Example response**

```json
[
    {
        "id": 1,
        "name": "nv",
        "aliases": [
            "non-verbal",
            "nonverbal"
        ],
        "description": "Hi! This person is nonverbal..",
        "created": "2021-01-19T18:05:43.147314Z"
    },
    // ...
]
```

### `GET /categories`

Gets all categories from the database. Returns an array of [category objects](#category-object) on success,
or `204 No Content` if there are no categories.

```
GET https://api.termora.org/v1/categories
```

**Example response**

```json
[
    {
        "id": 1,
        "name": "Plurality"
    },
    {
        "id": 2,
        "name": "LGBTQ+"
    }
]
```

### `GET /list`

Gets all terms from the database. Returns an array of [term objects](#term-object) on success,
or `204 No Content` if there are no terms.  
The query parameter `?flags=int` can be used to filter terms.
By default, a flag of `8` is used, which hides terms with the "hidden" flag.

**Example query**

```
GET https://api.termora.org/v1/list
```

**Example response**

```json
[
    {
        "id": 362,
        "category_id": 1,
        "category": "Plurality",
        "name": "Ability Divergence",
        "aliases": [
            "Ability Divergent"
        ],
        "description": "...",
        "source": "...",
        "created": "2021-04-02T18:06:13.37965Z",
        "last_modified": "2021-04-02T18:06:13.37965Z",
        "tags": [
            "Introjects/Introtives",
            "Plurality"
        ],
        "flags": 0
    },
    {
        "id": 234,
        "category_id": 2,
        "category": "LGBTQ+",
        "name": "Abro-",
        "aliases": [
            "Abrosexual",
            "Abroromantic"
        ],
        "description": "...",
        "source": "...",
        "created": "2021-03-19T19:49:54.449186Z",
        "last_modified": "2021-03-26T19:21:23.108998Z",
        "tags": [
            "LGBTQ+",
            "Romantic Orientation",
            "Sexuality"
        ],
        "flags": 0
    },
    // ...
]
```

### `GET /list/:id`

Gets all terms from the database. Returns an array of [term objects](#term-object) on success,
`204 No Content` if there are no terms, or `400 Bad Request` if `:id` was not an integer.

**Example query**

```
GET https://api.termora.org/v1/list/2
```

**Example response**

```json
[
    {
        "id": 234,
        "category_id": 2,
        "category": "LGBTQ+",
        "name": "Abro-",
        "aliases": [
            "Abrosexual",
            "Abroromantic"
        ],
        "description": "...",
        "source": "Unknown; already in circulation.",
        "created": "2021-03-19T19:49:54.449186Z",
        "last_modified": "2021-03-26T19:21:23.108998Z",
        "tags": [
            "LGBTQ+",
            "Romantic Orientation",
            "Sexuality"
        ],
        "flags": 0
    },
    {
        "id": 95,
        "category_id": 2,
        "category": "LGBTQ+",
        "name": "Agender",
        "aliases": [],
        "description": "...",
        "source": "https://lgbta.wikia.org/wiki/Agender",
        "created": "2021-02-02T17:10:58.830144Z",
        "last_modified": "2021-03-26T20:12:03.028932Z",
        "tags": [
            "LGBTQ+",
            "Gender",
            "Non-Binary"
        ],
        "flags": 0
    },
    // ...
]
```

## Version history

- **2021-04-08** (v1): initial documentation