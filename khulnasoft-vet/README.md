# khulnasoft-vet

[![Build Status](https://drone.khulnasoft.com/api/badges/khulnasoft/khulnasoft-vet/status.svg)](https://drone.khulnasoft.com/khulnasoft/khulnasoft-vet)

`go vet` tool for Khulnasoft

| Analyzer   | Description                                                                 |
|------------|-----------------------------------------------------------------------------|
| Imports    | Checks for import sorting. stdlib->github.com/khulnasoft->other                     |
| License    | Checks file headers for some form of `Copyright...YYYY...Khulnasoft/Gogs`        |
| Migrations | Checks for black-listed packages in `github.com/khulnasoft/khulnasoft/models/migrations` |
