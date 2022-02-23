import pyexiv2
import json


img = pyexiv2.Image(r'./regarding_geometry.png')
userdata="Congratulations you are very close indeed traveler. long it has been since you set out, but observant you have remained. Gateway:https://securesha.re/?#u=0d39c5a4b5eb0ef2dd919b468410b572.bin"

          
img.modify_comment(userdata)
print(img.read_comment())

