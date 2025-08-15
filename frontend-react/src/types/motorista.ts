export interface Motorista {
  id: string;
  nome: string;
  email: string;
  telefone?: string;
  status?: string;
  modelo_veiculo?: string;
  placa_veiculo?: string;
  criado_em?: string;
  documentos?: Documento[];
}

export interface Documento {
  id: string;
  tipo_documento: string;
  caminho_arquivo: string;
  formato: string;
  tamanho: number;
  status: string;
  criado_em: string;
}

export interface CadastroMotoristaPayload {
  nome: string;
  data_nascimento: string; // DD/MM/AAAA
  cpf: string;
  cnh: string;
  categoria_cnh: string;
  validade_cnh: string; // DD/MM/AAAA
  placa_veiculo: string;
  modelo_veiculo: string;
  telefone: string;
  email: string;
  senha: string;
  confirmacao_senha: string;
}
