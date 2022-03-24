from pydantic import BaseModel

# kelas ini sebagai jembatan antara aplikasi dg database
class Todo(BaseModel):
    name: str
    description: str
    completed: bool
