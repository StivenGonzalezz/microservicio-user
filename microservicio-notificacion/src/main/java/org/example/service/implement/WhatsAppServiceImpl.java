package org.example.service.implement;

import lombok.RequiredArgsConstructor;
import org.example.service.interfaces.WhatsAppService;
import org.springframework.stereotype.Service;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.*;
import org.springframework.web.client.RestTemplate;

import java.util.HashMap;
import java.util.Map;

@Service
@RequiredArgsConstructor
public class WhatsAppServiceImpl implements WhatsAppService {

    @Value("${whatsapp.api.url}")
    private String apiUrl;

    @Value("${whatsapp.phone.number.id}")
    private String phoneNumberId;

    @Value("${whatsapp.api.token}")
    private String accessToken;

    private final RestTemplate restTemplate = new RestTemplate();

    @Override
    public void sendMessage(String number, String message) {
        HttpHeaders headers = new HttpHeaders();
        headers.setContentType(MediaType.APPLICATION_JSON);
        headers.setBearerAuth(accessToken);

        // 🔹 Construcción del JSON para mensaje tipo "template" (SIN parámetros)
        Map<String, Object> template = new HashMap<>();
        template.put("name", "messagegeneric"); // 👈 nombre de la plantilla
        template.put("language", Map.of("code", "en_US")); // 👈 idioma correcto

        Map<String, Object> body = new HashMap<>();
        body.put("messaging_product", "whatsapp");
        body.put("to", number);
        body.put("type", "template");
        body.put("template", template);

        HttpEntity<Map<String, Object>> entity = new HttpEntity<>(body, headers);

        ResponseEntity<String> response = restTemplate.postForEntity(apiUrl, entity, String.class);

        System.out.println("Respuesta WhatsApp API: " + response.getBody());
    }
}
