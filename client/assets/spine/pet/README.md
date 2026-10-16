# Pet Spine assets

Place exported Spine files here after art production:

- `pet.json` or `pet.skel`
- `pet.atlas`
- `pet.png` (texture)

Animation names must match `docs/spine-pet.md`.

In Cocos: import this folder, add `sp.Skeleton` on PetRoot, assign skeleton data, attach `PetActor`.

Pixel art: use point filtering (nearest) on the texture.
