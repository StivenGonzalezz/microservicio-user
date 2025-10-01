package org.example.dto;
import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.NotBlank;

@JsonIgnoreProperties(ignoreUnknown = true)
public record CreateNotificationDTO(
        @NotBlank String affair,
        @NotBlank @Email String email,
        @NotBlank String body,
        @NotBlank String number
) {
}

