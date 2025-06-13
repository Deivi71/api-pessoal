USE apigestarbem;

-- Criar tabela de publicações
CREATE TABLE IF NOT EXISTS publicacoes(
    id int auto_increment primary key,
    titulo varchar(100) not null,
    conteudo text not null,
    autor_id int not null,
    curtidas int default 0,
    criadoEm timestamp default current_timestamp(),
    
    FOREIGN KEY (autor_id) 
    REFERENCES usuarios(id) 
    ON DELETE CASCADE
) ENGINE=INNODB;

-- Criar tabela de curtidas (para evitar curtidas duplicadas)
CREATE TABLE IF NOT EXISTS curtidas(
    usuario_id int not null,
    publicacao_id int not null,
    criadoEm timestamp default current_timestamp(),
    
    FOREIGN KEY (usuario_id) 
    REFERENCES usuarios(id) 
    ON DELETE CASCADE,
    
    FOREIGN KEY (publicacao_id) 
    REFERENCES publicacoes(id) 
    ON DELETE CASCADE,
    
    PRIMARY KEY (usuario_id, publicacao_id)
) ENGINE=INNODB;

-- Inserir dados de exemplo (opcional)
INSERT INTO publicacoes (titulo, conteudo, autor_id) VALUES 
('Minha primeira publicação', 'Este é um exemplo de publicação no Gestar Bem!', 1),
('Dicas de bem-estar', 'Compartilhando algumas dicas importantes para o dia a dia...', 1); 