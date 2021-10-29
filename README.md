Matilda is a waltzer.

The purpose of this project is to explore paths through floor. I am using
Location to represent a square on a grid. The square have transitions from
which we can build a graph. On initialization, we will build all paths from
one location to another. By calculating these ahead of time we can have
path tables which should make decision-making more efficient and hopefully
expose some other useful tricks such as collision prediction.
