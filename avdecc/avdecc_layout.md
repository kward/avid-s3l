# avdecc_layout

An avdecc_layout.xml file can be validated against the avdecc_layout.xsd using
the following steps.

1. Fetch the XSD file.

    ```
    curl -O https://raw.githubusercontent.com/kward/avid-venue/master/avdecc/avdecc_layout.xsd
    ```

2. Run `xmllint` using the local file.

    ```
    xmllint --noout --schema avdecc_layout.xsd avdecc_layout.xml
    ```

    
