import json
import matplotlib.pyplot as plt
from mpl_toolkits.mplot3d import Axes3D

# Step 1: Load the JSON data
with open('check/24_1734726292.json', 'r') as file:
    data = json.load(file)

# Step 2: Extract data for visualization
# Snake geometries
snake_x, snake_y, snake_z = [], [], []
for i, snake in enumerate(data['snakes']):
    for coord in snake['geometry']:
        snake_x.append(coord[0])
        snake_y.append(coord[1])
        snake_z.append(coord[2])

# Enemy geometries
enemy_x, enemy_y, enemy_z = [], [], []
for i, enemy in enumerate(data['enemies']):
    for coord in enemy['geometry']:
        enemy_x.append(coord[0])
        enemy_y.append(coord[1])
        enemy_z.append(coord[2])

# Food points
food_x, food_y, food_z = [], [], []
for food in data['food']:
    food_x.append(food['c'][0])
    food_y.append(food['c'][1])
    food_z.append(food['c'][2])

# Special food points (Golden and Suspicious)
special_x, special_y, special_z = [], [], []
for coord in data['specialFood']['golden']:
    special_x.append(coord[0])
    special_y.append(coord[1])
    special_z.append(coord[2])

for coord in data['specialFood']['suspicious']:
    special_x.append(coord[0])
    special_y.append(coord[1])
    special_z.append(coord[2])

# Fence coordinates (from `fences`)
fence_x, fence_y, fence_z = [], [], []
for coord in data['fences']:
    fence_x.append(coord[0])
    fence_y.append(coord[1])
    fence_z.append(coord[2])

# Step 3: Create the 3D plot
fig = plt.figure()
ax = fig.add_subplot(111, projection='3d')

# Step 4: Plot each type of entity
# ax.scatter(snake_x, snake_y, snake_z, c='blue', marker='o', label='Snakes')
ax.scatter(snake_x, snake_y, snake_z, c='blue', marker='o', label='Snakes')
ax.scatter(food_x, food_y, food_z, c='green', marker='s', label='Food')
ax.scatter(special_x, special_y, special_z, c='gold', marker='x', label='Special Food')
ax.scatter(fence_x, fence_y, fence_z, c='black', marker='d', label='Fences')

# Step 5: Customize the plot
ax.set_xlabel('X Axis')
ax.set_ylabel('Y Axis')
ax.set_zlabel('Z Axis (Type)')

plt.title('3D Visualization of Snake Game Data')
ax.legend()
plt.show()
