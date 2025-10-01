package org.example.service.implement;

import lombok.RequiredArgsConstructor;
import org.example.dto.EmailDTO;
import org.example.service.interfaces.EmailService;
import org.simplejavamail.api.email.Email;
import org.simplejavamail.api.mailer.Mailer;
import org.simplejavamail.api.mailer.config.TransportStrategy;
import org.simplejavamail.email.EmailBuilder;
import org.simplejavamail.mailer.MailerBuilder;
import org.springframework.scheduling.annotation.Async;
import org.springframework.stereotype.Service;

@Service
@RequiredArgsConstructor
public class EmailServiceImpl implements EmailService {


    String host = "smtp.gmail.com";
    String username = "proyectofinal255@gmail.com";
    int port = 587;
    String password = "eiqawevrvhabdhmn";

    @Override
    @Async
    public void sendEmail(EmailDTO emailDTO) throws Exception {

        Email email = EmailBuilder.startingBlank()
                .from(username)
                .to(emailDTO.destinatario())
                .withSubject(emailDTO.asunto())
                .withHTMLText(emailDTO.cuerpo())
                .buildEmail();


        try (Mailer mailer = MailerBuilder
                .withSMTPServer(host, port, username, password)
                .withTransportStrategy(TransportStrategy.SMTP_TLS)
                .withDebugLogging(true)
                .buildMailer()) {


            mailer.sendMail(email);
        }

    }
    @Override
    @Async
    public void sendEmailAttach(EmailDTO emailDTO, byte[] contenidoPDF, String nombreArchivo) throws Exception {
        Email email = EmailBuilder.startingBlank()
                .from(username)
                .to(emailDTO.destinatario())
                .withSubject(emailDTO.asunto())
                .withHTMLText(emailDTO.cuerpo())
                .withAttachment(nombreArchivo, contenidoPDF, "application/pdf")
                .buildEmail();

        try (Mailer mailer = MailerBuilder
                .withSMTPServer(host, port, username, password)
                .withTransportStrategy(TransportStrategy.SMTP_TLS)
                .withDebugLogging(true)
                .buildMailer()) {

            mailer.sendMail(email);
        }
    }
}
