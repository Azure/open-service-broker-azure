# OSBA PCF Tile Generating

## Prerequisites

  * make
  * Docker (version 17.04 or greater)
  * Docker Compose (version 1.16.1 or greater)

## Generate a PCF Tile

  1. Under the folder `pcf-tile`, execute below command.
  
  * Build a major version:
    ```
    make generate-major-tile
    ```
    
  * Build a minor version:
    ```
    make generate-minor-tile
    ```
    
  * Build a patch version:
    ```
    make generate-tile
    ```

  2. You can find open-service-broker-azure-`VERSION`.pivotal under the folder `product`.
