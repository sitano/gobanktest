# gobanktest

Example of simple server for bank transactions processing.

Component schema:

                     +----------+            +----------------+
                     | Client   +------------+                |
                     |          |            | API            |
                     +----------+            | * /transaction |
                                             | * /balances    |
                              +----------+   |                |
                              | Server   +---+                |
                              |          |   -----------------+
                              |          |            |        
                              +----------+            |        
    +--------------+                                  |        
    | Domain       |   +-----------------+   +--------v-------+
    |              |   | Storage         |   | Business Logic |
    |   +--------+ |   | * save, load    |   | * /transaction |
    |   |User    | |   | * mt safe       |   | * /balances    |
    |   |  +-----+ +---+ * tx support    <--->                |
    |   |  |Purse| |   |   * chg balance |   |                |
    |   |  |     | |   |                 |   |                |
    |   +--------+ |   |                 |   |                |
    |              |   |                 |   |                |
    +--------------+   +-----------------+   +----------------+
