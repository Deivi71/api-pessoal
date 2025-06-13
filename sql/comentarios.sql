USE apigestarbem;

-- Criar tabela de comentários
CREATE TABLE IF NOT EXISTS comentarios(
    id int auto_increment primary key,
    conteudo text not null,
    autor_id int not null,
    publicacao_id int not null,
    criadoEm timestamp default current_timestamp(),
    
    FOREIGN KEY (autor_id) 
    REFERENCES usuarios(id) 
    ON DELETE CASCADE,
    
    FOREIGN KEY (publicacao_id) 
    REFERENCES publicacoes(id) 
    ON DELETE CASCADE
) ENGINE=INNODB;

-- Inserir dados de exemplo (opcional)
INSERT INTO comentarios (conteudo, autor_id, publicacao_id) VALUES 
('Ótima publicação! Muito útil.', 2, 1),
('Concordo totalmente com você!', 3, 1),
('Obrigado por compartilhar essas dicas.', 2, 2); 