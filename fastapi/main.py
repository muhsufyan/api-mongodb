from fastapi import FastAPI
from routes.todos_route import todo_api_router
from fastapi.responses import ORJSONResponse
app = FastAPI()
# masukkan semua route di direktory routes
app.include_router(todo_api_router)