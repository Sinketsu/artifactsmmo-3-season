package gen

//go:generate wget -O openapi.json https://api-test.artifactsmmo.com/openapi.json
//go:generate ogen -target=oas -package oas -clean openapi.json
