import hashlib
import threading
import os

NUM_THREADS = 4

x = "printdata.mcap"
fileSize=os.path.getsize(x)
bytesPerThread = fileSize / NUM_THREADS

def getSHA(filename, startIndex, hashChunkList):

    fileHash = hashlib.new("sha256")

    if startIndex < 0:
        startIndex = 0
        
    with open(filename, "rb") as file:
        i=0
        while i < (bytesPerThread / 4096):
            print("hashing chunk " + str(i), end='\r')
            chunk = file.read(4096)
            if not chunk:
                break
            fileHash.update(chunk)
            i += 1
            
    print("hashing complete.")
    return fileHash.hexdigest()
    
    
    
hl = []
threads = []

for i in range(NUM_THREADS):
    threads.append(threading.Thread(target=getSHA, args=(x, (bytesPerThread * i) - 1, hl), daemon=True))
    threads[i].start

for i in range(NUM_THREADS):
    hl.append(threads[i].join)

print(hl)
