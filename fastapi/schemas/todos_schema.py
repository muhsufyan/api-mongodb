# this file used catch data (filter request & filter response)
# catch single todo data. returnnya adlh dict
def todo_serializer(todo):
    return {
        "id": str(todo["_id"]),
        "name": todo["name"],
        "description": todo["description"],
        "completed": todo["completed"]
    }

# catch more todos data. returnnya adlh list
def todos_serializer(todos):
    return [todo_serializer(todo) for todo in todos]
