module ex6/server

go 1.17

replace ex6/fileservice => ../fileservice

require google.golang.org/grpc v1.45.0
require ex6/fileservice v0.0.0-00010101000000-000000000000