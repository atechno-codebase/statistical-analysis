l=int(input("length: "))
data=[float(input(str(i)+ ": ")) for i in range(l)]
print(data)

mean=sum(data)/l

# stddev=(sum([(x-mean)**2 for x in data])/(l-1))**0.5
# print("stddev: ", stddev)
stddev=(sum([(x-mean)**2 for x in data])/(l))**0.5
print("stddev: ", stddev)
