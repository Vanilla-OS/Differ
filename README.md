# Differ

Differ is a REST-API for retrieving package changes between images on immutable Linux distributions.
It was designed for allowing the user to visualize changes in package versions in [Vanilla OS](https://vanillaos.org), but this tool can be used by any distribution without any changes to the code.

## Build and Setup

In order to setup Differ, you need Go 1.21 and some Sqlite manager (in this section we use the `sqlite3` command line tool).
First, you need to create a database to store all the images and releases Differ will manage, as well as auth information (see read-only subsection below).
Using the `sqlite3` tool, run the following commands:

```sh
$ sqlite3 /path/to/database.db
sqlite> create table auth("ID" INTEGER, name, pass TEXT, PRIMARY KEY("ID"));
sqlite> insert into auth values(1, "admin_user", "admin_password"); # Replace user and password with something secure
```

If everything was set up correctly, the auth table should look like this:

```sh
sqlite> select * from auth;
1|admin_user|admin_password
```

Now, build the project with the provided Makefile and run Differ passing the database path as argument.

```sh
$ make
$ ./differ path/to/database.db
```

When setting up Differ in production, you must `export GIN_MODE=release` before running the binary.

## Endpoints

### Status

Simple check to see if the server is running correctly.

*Endpoint:* `http://[base_url]/status`

*Parameters:* None

*Returns:*

- `200 OK` on success.

``` json
{
    "status": "ok"
}
```

### Images

Images are, as the name implies, image types shipped by the distribution.
For example, a distro can ship both GNOME and a KDE versions as separate images, which would can be managed with the endpoints below:

#### Create image

Creates a new image in the dabatase. Every release (see subsection below) is attached to an image.

*Endpoint:* `http://[base_url]/images/new`

*Parameters:*

- *Name:* Image name
- *URL:* Where the image is hosted or its repository. For information purposes only.

```json
{
    "name": "pico",
    "url": "https://github.com/Vanilla-OS/pico-image"
}
```

*Returns:*

- `200 OK` on success.

#### Get image by name

Retrieves information about an image given its name.

**Endpoint:* `http://[base_url]/images/[name]`

*Parameters:* None

*Returns:*

- `200 OK` on success, alongside the image information.

``` json
{
    "image": {
        "name": "pico",
        "url": "https://github.com/Vanilla-OS/pico-image",
        "releases": [
            ...
        ]
    }
}
```

- `400 Bad Request` if image cannot be found.

#### Get all images

Retrieves information about all images.

*Endpoint:* `http://[base_url]/images`

*Parameters:* None

*Returns:*

- `200 OK` on success, alongside the images information.

``` json
{
    "images": [
        {
            "name": "pico",
            "url": "https://github.com/Vanilla-OS/pico-image",
            "releases": [
                ...
            ]
        },
        {
            "name": "core",
            "url": "https://github.com/Vanilla-OS/core-image",
            "releases": [
                ...
            ]
        }
    ]
}
```

### Releases

A release is a new version of some image, where packages can be added, removed, upgraded or downgraded. A release also has a digest, which should be the image digest provided by the container manager.

#### Create release

Creates a new release for the given image.

*Endpoint:*`http://[base_url]/images/[image]/new`

*Parameters:*

- *Digest:* Image digest
- *Packages:* List of package names and versions in the current release. `get_all_packages.py` contains a script for extracting the list from a Debian-based distribution.

```json
{
    "digest": "sha256:a99e4593b23fd07e3761639e9db38c0315e198d6e39dad6070e0e0e88be3de0d",
    "packages": [
        {
            "name": "apt",
            "version": "2.7.3"
        },
        ...
    ]
}
```

*Returns:*

- `200 OK` on success.

#### Get latest release for image

Retrieves the most recent release from the image.

*Endpoint:* `http://[base_url]/[image]/latest`

*Parameters:* None

*Returns:*

- `200 OK` on success, alongside the release information.

``` json
{
    "release": {
        "digest": "sha256:a99e4593b23fd07e3761639e9db38c0315e198d6e39dad6070e0e0e88be3de0d",
        "date": "2023-11-26T12:56:53.377838864Z",
        "packages": [
            ...
        ]
    }
}
```

#### Get specific release

Searches for a specific release by its digest.

*Endpoint:* `http://[base_url]/[image]/[digest]`

*Parameters:* None

*Returns:*

- `200 OK` on success, alongside the release information.

``` json
{
    "release": {
        "digest": "sha256:a99e4593b23fd07e3761639e9db38c0315e198d6e39dad6070e0e0e88be3de0d",
        "date": "2023-11-26T12:56:53.377838864Z",
        "packages": [
            ...
        ]
    }
}
```

- `400 Bad Request` if release cannot be found.

#### Release diff

The most important endpoint in the API. Given two digests, generates a list of changed packages. This information is cached so future queries are nearly instant.

*Endpoint:* `http://[base_url]/images/[image]/diff`

*Parameters:*

- *Old digest:* Digest of the older image, which is usually the image the user is currently on.
- *New digest:* Digest of the newer image, which is usually the image the user wants to update to.

```json
{
    "old_digest": "sha256:a99e4593b23fd07e3761639e9db38c0315e198d6e39dad6070e0e0e88be3de0c",
    "new_digest": "sha256:a99e4593b23fd07e3761639e9db38c0315e198d6e39dad6070e0e0e88be3de0d"
}
```

*Returns:*

- `200 OK` on success, alongside the modified packages.

``` json
{
    "_new_digest": "sha256:a99e4593b23fd07e3761639e9db38c0315e198d6e39dad6070e0e0e88be3de0d",
    "_old_digest": "sha256:a99e4593b23fd07e3761639e9db38c0315e198d6e39dad6070e0e0e88be3de0c",
    "added": [
        {
            "name": "zsh-common",
            "new_version": "5.9-5"
        }
    ],
    "downgraded": [
        {
            "name": "autoconf",
            "old_version": "2.71-3"
        }
    ],
    "removed": [
        {
            "name": "nano",
            "old_version": "6.2"
        }
    ],
    "upgraded": [
        {
            "name": "curl",
            "old_version": "1.0.3-aplha",
            "new_version": "8.2.1-2"
        }
    ]
}
```

- `400 Bad Request` if either digest cannot be found.
