# Dabbot [![Docker Image Build](https://github.com/Bermos/Dabbot/actions/workflows/docker-image.yml/badge.svg)](https://github.com/Bermos/Dabbot/actions/workflows/docker-image.yml)

## Why
Did you really expect a reason for why one would need a dabbot?

## How-To
### Quickstart
1. Download/Clone repository
2. Run `go mod download` to download dependencies
3. Run `go build -o dabbot .` to compile
4. Create `dabs` folder and fill with the appropriate files.
5. Set environmental variable `TOKEN` for the bot to use
6. Start the bot with `./dabbot` or do step 4. & 5. with `TOKEN=abc:def ./dabbot`

### Docker
1. Download/Clone repository
2. Create `dabs` folder and fill with the appropriate files.
3. Build docker image and run it `docker build -t dabbot . && docker run --rm -e TOKEN='abc:def' -v dabs:/app/dabs dabbot`

### Contribute
1. Fork it
2. Clone it: `git clone https://github.com/bermos/Dabbot`
3. Create your feature branch: `git checkout -b my-new-feature`
4. Make changes and add them: `git add .`
5. Commit: `git commit -m 'Add some feature'`
6. Push: `git push origin my-new-feature`
7. Pull request
