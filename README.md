# Scrobble Fetcher in Go

This is a Scrobble fetcher, written in Go (programming language). It fetches scrobbles from Last.fm, then shows the output in your terminal.

## Usage

**Clone repository**

```sh
mkdir -p /$HOME/go-projects/go-fetch-scrobbles && cd $_
```

```sh
git clone https://github.com/ricardobalk/go-fetch-scrobbles .
```

**Run app**

- To receive a list of [Batmaniosaurus'](https://last.fm/user/Batmaniosaurus) 50 last played tracks:

  ```sh
  go run main.go --api-token "$LAST_FM_API_TOKEN" --username "Batmaniosaurus"
  ```

<details><summary>Example output</summary><pre>01: Nina June - Shadows & Riddles [Shadows & Riddles] at 2020-07-03 14:12:40 +0000 UTC
02: Miss Montreal - House Upon The Hill [House Upon The Hill] at 2020-07-03 14:09:23 +0000 UTC
03: Angela Moyra - You, Me & the Sea [You, Me & the Sea] at 2020-07-03 14:05:16 +0000 UTC
04: Racoon - Het Is Al Laat Toch [Het Is Al Laat Toch] at 2020-07-03 14:01:47 +0000 UTC
05: Nina June - Summersnow [Summersnow] at 2020-07-03 13:58:25 +0000 UTC
...
50: Bløf - Alles is Liefde [Hier - Het Beste Van 20 Jaar Bløf] at 2020-07-01 18:42:13 +0000 UTC</pre</details>

- To receive JSON-formatted output:

```sh
go run main.go --api-token "$LAST_FM_API_TOKEN" --username "Batmaniosaurus" --format "json"
```

<details><summary>Example output</summary>
<pre>[{"artist":"Nina June","song":"Shadows \u0026 Riddles","album":"Shadows \u0026 Riddles","when":{"posix":1593785703,"human":"2020-07-03 14:15:03 +0000 UTC"}},{"artist":"Miss Montreal","song":"House Upon The Hill","album":"House Upon The Hill","when":{"posix":1593785363,"human":"2020-07-03 14:09:23 +0000 UTC"}},{"artist":"Angela Moyra","song":"You, Me \u0026 the Sea","album":"You, Me \u0026 the Sea","when":{"posix":1593785116,"human":"2020-07-03 14:05:16 +0000 UTC"}},{"artist":"Racoon","song":"Het Is Al Laat Toch","album":"Het Is Al Laat Toch","when":{"posix":1593784907,"human":"2020-07-03 14:01:47 +0000 UTC"}}, ...</pre>
</details>

- To receive raw response data from Last.fm:
```sh
go run main.go --api-token "$LAST_FM_API_TOKEN" --username "Batmaniosaurus" --format "raw"
```