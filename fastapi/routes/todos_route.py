from fastapi import APIRouter
from config.database import collection_name
from models.todos_model import Todo
from schemas.todos_schema import todo_serializer, todos_serializer
from bson import ObjectId
# custom filter response
from fastapi.responses import ORJSONResponse

todo_api_router =  APIRouter(tags=["todo"])

# get all data
@todo_api_router.get("/")
async def get_todos():
    # select * from tabel;
    dataset = todos_serializer(collection_name.find())
    return{"status": "ok","data": dataset}

# get one data
@todo_api_router.get("/{id}")
async def get_todo(id: str):
    data = todos_serializer(collection_name.find({"_id":ObjectId(id)}))
    return {"status":"ok","data":data}

# buat data todo baru
@todo_api_router.post("/")
async def create_todo(todo: Todo):
    _id = collection_name.insert_one(dict(todo))
    todo = todos_serializer(collection_name.find({"_id": _id.inserted_id}))
    return {"status":"ok","data":todo}

# update
@todo_api_router.put("/{id}")
async def update_todo(id: str, todo: Todo):
    # cari id data todo yg diupdate lalu update datanya
    collection_name.find_one_and_update({"_id": ObjectId(id)}, {
        "$set": dict(todo)
    })
    todo = todos_serializer(collection_name.find({"_id": ObjectId(id)}))
    return {"status":"ok","data":todo}

# delete
@todo_api_router.delete("/{id}")
async def delete_todo(id: str):
    # cari id data todo yg diingin dihapus lalu hapus datanya
    collection_name.find_one_and_delete({"_id": ObjectId(id)})
    return {"status": "ok", "data":[]}
