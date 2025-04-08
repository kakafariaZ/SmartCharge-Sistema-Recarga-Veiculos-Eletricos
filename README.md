# SmartCharge 🚗⚡

Sistema inteligente para gerenciamento e reserva de pontos de recarga para veículos elétricos.

---

## 📖 Descrição do Projeto

O **SmartCharge** é uma solução inovadora para facilitar o uso de veículos elétricos, oferecendo funcionalidades como:
- Localização e disponibilidade de pontos de recarga em tempo real.
- Filas para gerenciamento das recagas.
- Distribuição inteligente da demanda para reduzir o tempo de espera.
- Liberação automática de carregadores após a recarga.
- Pagamento integrado.

O sistema utiliza comunicação via TCP para troca de mensagens entre os carros, postos de recarga e o servidor central, garantindo eficiência e escalabilidade.

---

## 🚀 Como executar o sistema

### Pré-requisitos
- Docker e Docker Compose instalados na máquina.
- Porta `8080` disponível para o servidor.

### Passos para execução

1. **Clonar o repositório**
   ```bash
   git clone https://github.com/SeuUsuario/SmartCharge-Sistema-Recarga-Veiculos-Eletricos.git
   cd SmartCharge-Sistema-Recarga-Veiculos-Eletricos
   ```

2. **Construir e iniciar os containers**
   Execute o comando abaixo para construir as imagens e iniciar os containers:
   ```bash
   docker-compose up --build
   ```

3. **Escalar o número de carros (opcional)**
   Para criar múltiplas instâncias do cliente `client-car`, use o comando:
   ```bash
   docker-compose up --scale client-car=3
   ```
   O número `3` pode ser ajustado conforme necessário.

4. **Parar os containers**
   Para parar e remover os containers, execute:
   ```bash
   docker-compose down
   ```

---

## 🧪 Scripts de Experimentos

### 1. **Executar o servidor separadamente**
   Para iniciar apenas o servidor:
   ```bash
   docker-compose up server
   ```

### 2. **Executar um posto de recarga**
   Para iniciar um posto de recarga:
   ```bash
   docker-compose up client-station
   ```

### 3. **Executar um carro**
   Para iniciar um carro individualmente:
   ```bash
   docker-compose run --service-ports --name car1 client-car
   ```

   Para iniciar múltiplos carros, altere o nome (`car1`, `car2`, etc.) e repita o comando.

---

## 📂 Estrutura do Projeto

- **[`server`](server )**: Contém o código do servidor principal que gerencia as conexões e distribui os postos de recarga.
- **[`client-car`](client-car )**: Código do cliente que simula os carros elétricos.
- **[`client-station`](client-station )**: Código do cliente que simula os postos de recarga.
- **[`Internal`](Internal )**: Contém módulos internos, como gerenciamento de filas e comunicação via WebSocket.
- **[`models`](models )**: Definições de modelos de dados, como `ChargeStation`.

---

## 📝 Notas Importantes

- O sistema utiliza comunicação via TCP para troca de mensagens entre os carros, postos e o servidor.
- Certifique-se de que o arquivo docker-compose.yml está configurado corretamente para o ambiente local.
- Para simular diferentes cenários, ajuste os parâmetros no código ou no arquivo docker-compose.yml.

---

## ⚙️ Scripts Automatizados

O projeto inclui o arquivo docker-compose.yml para facilitar a execução e o gerenciamento dos containers Docker. Com ele, é possível:
- Construir e iniciar todos os serviços com um único comando (`docker-compose up`).
- Escalar o número de instâncias de carros ou postos de recarga.
- Parar e remover todos os containers com facilidade (`docker-compose down`).

---

## 📧 Autores

- Kauan Caio de Arruda Farias `fariak@ecomp.uefs.br`.
- Nathielle Cerqueira Alves `nalves@ecomp.uefs.br`.
- Vitória Tanan dos Santos `vtsantos@ecomp.uefs.br`.
