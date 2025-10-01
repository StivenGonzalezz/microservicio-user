package org.example.service.interfaces;

import org.example.dto.EmailDTO;

public interface EmailService {
    void sendEmail(EmailDTO emailDTO) throws Exception;
    void sendEmailAttach(EmailDTO emailDTO, byte[] contenidoPDF, String nombreArchivo) throws Exception;

}
