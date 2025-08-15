export function formatDate(value?: string | Date): string {
  if (!value) return '-';
  const d = typeof value === 'string' ? new Date(value) : value;
  if (isNaN(d.getTime())) return '-';
  return d.toLocaleDateString('pt-BR');
}

export function formatCPF(cpf?: string): string {
  if(!cpf) return '-';
  const digits = cpf.replace(/\D/g,'');
  if(digits.length !== 11) return cpf;
  return digits.replace(/(\d{3})(\d{3})(\d{3})(\d{2})/, '$1.$2.$3-$4');
}

export function formatTelefone(t?: string): string {
  if(!t) return '-';
  const d = t.replace(/\D/g,'');
  if(d.length===11) return d.replace(/(\d{2})(\d{5})(\d{4})/, '($1) $2-$3');
  if(d.length===10) return d.replace(/(\d{2})(\d{4})(\d{4})/, '($1) $2-$3');
  return t;
}

export function formatPlaca(p?: string): string {
    if (!p) return '-';
    return p.toUpperCase();
}

// Sanitizers (enviar sempre limpo ao backend)
export const unmask = (v: string) => v.replace(/\D/g,'');
export const sanitizeCPF = (v: string) => unmask(v).slice(0,11);
export const sanitizeTelefone = (v: string) => unmask(v).slice(0,11);
export const sanitizePlaca = (v: string) => v.toUpperCase().replace(/[^A-Z0-9]/g,'').slice(0,7);
export const sanitizeCNH = (v: string) => unmask(v).slice(0,11);
export const sanitizeEmail = (v: string) => v.trim();
