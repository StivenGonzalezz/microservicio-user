package org.example.dto;
import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.NotBlank;

public record EmailDTO (
    @NotBlank  String asunto,
    @NotBlank  String cuerpo,
    @NotBlank @Email String destinatario)
{
}
