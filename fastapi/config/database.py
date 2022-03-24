from pymongo import MongoClient
import pymongo

# konek ke db
# client = pymongo.MongoClient("mongodb://admin:password@localhost:27017/todo_application?authSource=admin&authMechanism=SCRAM-SHA-256")
client = pymongo.MongoClient("mongodb://localhost:27017/todo_application")

# create database todo_application
db = client["todo_application"]
# db = client.todo_application
# pd mongodb setiap data disimpan dlm collection (tabel)
# buat collection dg nama todos_app
collection_name = db["todos_app"]
print(collection_name)