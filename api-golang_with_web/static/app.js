async function buscarPessoa() {
  const termo = document.getElementById('busca').value.trim().toLowerCase();
  const res = await fetch('/api/pessoas');
  const pessoas = await res.json();

  const resultados = pessoas.filter(p =>
    p.nome.toLowerCase().includes(termo) || p.cpf.includes(termo)
  );

  const container = document.getElementById('resultadoBusca');
  container.innerHTML = '';

  if (resultados.length > 0) {
    resultados.forEach(resultado => {
      const div = document.createElement('div');
      div.textContent = `${resultado.nome} (${resultado.cpf}) - ${resultado.nascimento}`;

      const btn = document.createElement('button');
      btn.textContent = 'Excluir';
      btn.onclick = async () => {
        await fetch('/api/pessoas/' + resultado.id, { method: 'DELETE' });
        div.remove();
      };

      div.appendChild(btn);
      container.appendChild(div);
    });
  } else {
    container.textContent = 'Nenhuma pessoa encontrada.';
  }
}

function limparBusca() {
  document.getElementById('busca').value = '';
  document.getElementById('resultadoBusca').innerHTML = '';
}

document.getElementById('formPessoa').addEventListener('submit', async (e) => {
  e.preventDefault();
  const nome = document.getElementById('nome').value;
  const cpf = document.getElementById('cpf').value;
  const nascimento = document.getElementById('nascimento').value;

  const res = await fetch('/api/pessoas', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ nome, cpf, nascimento })
  });

  if (res.ok) {
    alert("Pessoa cadastrada!");
    e.target.reset();
  } else {
    alert("Erro ao cadastrar.");
  }
});
