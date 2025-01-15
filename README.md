# API ADV SYSTEM

#### Importanto o Projeto

Para evitar problemas no projeto dentro do VS Code, segue este padrão de importação.

Para isso é necessário abrir o arquivo `api-adv-system.code-workspace`

#### Start
Para iniciar todos os serviços dependente do projeto, execute o seguinte comando dentro do diretorio principal do projeto:

```sh
docker-compose -f dev/docker/docker-compose.yaml up --build
```

---

#### Banco de Dados
Para realizar uma conexão através do client de preferencia, basta utilizar as credenciais abaixo:

- **host**: 0.0.0.0
- **user**: root
- **password**: root
- **database**: adv_system

---  

#### Swagger 

Para executar o swagger/api localmente, no VS Code é necessário executar no "Run and Debug" da ferramenta.

O swagger em LocalHost esta na seguinte URL: 

➡️ [API Checkout](http://localhost:9000/docs/swagger/index.html#/)
