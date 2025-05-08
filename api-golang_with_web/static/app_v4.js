async function fetchPessoas() {
    const res = await fetch('/api/pessoas');
    const pessoas = await res.json();
    const lista = document.getElementById('listaPessoas');
    lista.innerHTML = '';
    pessoas.forEach(p => {
      const item = document.createElement('li');
      item.textContent = `${p.nome} (${p.cpf}) - ${p.nascimento}`;
      const del = document.createElement('button');
      del.textContent = 'Excluir';
      del.onclick = () => deletePessoa(p.cpf);
      item.appendChild(del);
      lista.appendChild(item);
    });
  }
  
  async function deletePessoa(cpf) {
    await fetch(`/api/pessoas/${cpf}`, { method: 'DELETE' });
    fetchPessoas();
  }
  
  document.getElementById('formPessoa').addEventListener('submit', async (e) => {
    e.preventDefault();
    const nome = document.getElementById('nome').value;
    const cpf = document.getElementById('cpf').value;
    const nascimento = document.getElementById('nascimento').value;
  
    await fetch('/api/pessoas', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ nome, cpf, nascimento })
    });
  
    fetchPessoas();
    e.target.reset();
  });
  
  fetchPessoas();
  