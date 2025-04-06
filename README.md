### Chaiwala Backend

##### TODO

- check why some userid is pgtype vs int32
- moved routes likely have extra json field requirements that should now come from path
- add json validations
- wipe old refresh token
- copy models from db and clean up the json serialization
- handle db errors better on key constraints should honestly make that default
- logging
- make all list endpoints paginated
  - add offset and page_size to query params
- find way to remove passwordDigestoffset
