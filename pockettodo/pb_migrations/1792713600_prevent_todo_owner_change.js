/// <reference path="../pb_data/types.d.ts" />
migrate((app) => {
  const collection = app.findCollectionByNameOrId("pbc_2097402128")
  collection.updateRule = '@request.auth.id != "" && user = @request.auth.id && @request.body.user:changed = false'
  return app.save(collection)
}, (app) => {
  const collection = app.findCollectionByNameOrId("pbc_2097402128")
  collection.updateRule = '@request.auth.id != "" && user = @request.auth.id'
  return app.save(collection)
})
