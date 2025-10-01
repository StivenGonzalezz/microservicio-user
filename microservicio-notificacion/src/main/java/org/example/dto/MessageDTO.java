package org.example.dto;

public record MessageDTO<T>(
        boolean error,
        T respuesta
) {
}