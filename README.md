### Chaiwala Backend

##### TODO

- check why some userid is pgtype vs int32
- moved routes likely have extra json field requirements that should now come from path
- add json validations
  - lets just use the go playground validator lib
- copy models from db and clean up the json serialization
- handle db errors better on key constraints should honestly make that default
- logging
  - need to actually add the logs
- make all list endpoints paginated
  - add offset and page_size to query params
- keep ingredients list, call it required ingredients or something as just list[str], nothing in DB for now
- add tea types
- im actually not a big fan of these logging. for the context key. lets remove and just make them constants elsewhere
- revoke refresh on Backend
- user should be able to see their recipes, whether it is public or not, but others should only be able to see public recipes. pattern will be stopped on frontend, but backend should have as well.

##### Done from TODO:

- implement recipe steps
- auto refresh on frontend

- find way to remove passwordDigestoffset
  - we are just gonna wipe the json, even thought ur not supposed to update the file
- remove some brew time stuff, can add later if necessary for filtering purposes
