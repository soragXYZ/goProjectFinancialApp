# Documentation

The documentation you are currently reading was made using [Material for MkDocs](https://squidfunk.github.io/mkdocs-material/getting-started/){:target="_blank"}.

Material for MkDocs is a powerful documentation framework on top of [MkDocs](https://www.mkdocs.org/){:target="_blank"}, a static site generator for project documentation.  

Documentation source files are written in Markdown, and configured with a single YAML configuration file.

## How is this documentation organized ?
For information, the documentation files are in the [docs](https://github.com/soragXYZ/freenahi/tree/main/docs){:target="_blank"} folder, located at root level of the repository.
```shell
freenahi
├── backend
│  ├── ...
│  └── ...
│
├── frontend
│  ├── ...
│  └── ...
|
├── docs
│  ├── ...
│  └── ...
|
├── .gitignore
├── LICENCE
└── README.md

```

Inside the **docs** folder, you will find the current documentation structure, which is quite simple
```shell
docs
├── docs
│  ├── assets
|  |  ├── ...
│  │  └── ...
│  ├── en
|  |  ├── ...
│  │  └── ...
│  └── fr
|     ├── ...
│     └── ...
│
├── overrides
│  ├── home-en.html
│  ├── home-fr.html
│  └── main.html
│
├── mkdocs.yaml
└── requirements.txt
```

## Build the documentation locally
If you want to contribute to this documentation, you can do it in 2 ways :

* Manually edit a file with Github directly, with the "Edit this page" button of the top right corner of this page.  
This can be helpful for correcting typos really fast
* Build the documentation locally to write news pages

???+ info
    Building the documentation locally should be the standard way

???+ note

    Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla et euismod
    nulla. Curabitur feugiat, tortor non consequat finibus, justo purus auctor
    massa, nec semper lorem quam in massa.

### Install the dependencies
First, you need to have Python installed on your system.
At the time of writing (April 2025), the last version is Python 3.13

=== "Windows"

    Download and install the Python MSI from [the official download page](https://www.python.org/downloads/){:target="_blank"}

=== "Linux"

    Python should be avalaible with your system package tool such as apt or dnf

    ```shell
    sudo apt install python3.13
    ```

### Create a virtual environment
A virtual environment is a folder that you pip, the packet manager of Python, install modules into. Every project should have its own virtual environment to avoid the problem where two different projects need two different versions of the same module in order to both work properly.

Navigate into the **docs** folder, then create the virtual environment here

=== "Windows"
    ```powershell
    py -m venv <NameOfTheVenv>
    ```

=== "Linux"
    ```shell
    python3.13 -m venv <NameOfTheVenv>
    ```

### Activate the virtual environment
=== "Windows"
    ```powershell
    <NameOfTheVenv>/Scripts/activate
    ```

=== "Linux"
    ```shell
    source <NameOfTheVenv>/bin/activate
    ```

### Install the dependencies
```shell
pip install -r requirements.txt
```

### Build and serve the doc
```shell
mkdocs serve
```

???+ warning
    You may have some warnings related to the git-committers plugin (API token missing or 429 HTTP errors).  
    This plugin is used to display the name of contributors associated with a file.

    You can follow the [documentation and create a Github token](https://github.com/byrnereese/mkdocs-git-committers-plugin?tab=readme-ov-file#config){:target="_blank"}, or simply deactivate the plugin in the mkdocs.yml file (enabled: false)

### Access the documentation
Your documentation should be up and running. Head over to [http://127.0.0.1:8000/freenahi](http://127.0.0.1:8000/freenahi){:target="_blank"} to see your documentation !