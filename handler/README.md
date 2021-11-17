# handler

`handler` implement `conn.Handler` interface. read data from reader and convert it to handleable data struct and call for execute it. Finally return response bytes for write back to clients.